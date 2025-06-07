package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/handlers/wishlist"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	cfg     *config.Config
	storage *storage.Storage
}

func New(cfg *config.Config, storage *storage.Storage) *Server {
	return &Server{cfg, storage}
}

func (s *Server) Start() {
	router := s.createRouter()

	srv := &http.Server{
		Addr:    s.cfg.App.Address,
		Handler: router,
	}

	serverErr := make(chan error, 1)

	go func() {
		fmt.Printf("Service Wishlister starting on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
		close(serverErr)
	}()

	quit := make(chan os.Signal, 1)

	select {
	case <-quit:
		fmt.Println("\nShutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Printf("Server forced to shutdown: %v\n", err)
		}
		fmt.Println("Server exited properly")
	case err := <-serverErr:
		fmt.Printf("Server error: %v\n", err)
	}
}

func (s *Server) createRouter() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(corsMiddleware)

	wishlisthandler := wishlist.New(s.cfg, s.storage, slog.Default())

	router.Get("/{user_id}", wishlisthandler.GetWishlist)
	router.Get("/{wish_id}/{user_id}", wishlisthandler.GetWish)
	router.Post("/", wishlisthandler.AddToWishlist)
	router.Delete("/{wish_id}/{user_id}", wishlisthandler.RemoveFromWishlist)
	router.Patch("/{wish_id}/{user_id}", wishlisthandler.PartialUpdateWish)

	return router
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
