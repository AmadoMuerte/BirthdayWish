package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/config"
	authhandler "github.com/AmadoMuerte/BirthdayWish/API/internal/http_server/handlers/auth"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/http_server/handlers/wishlist"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

type Server struct {
	cfg       *config.Config
	storage   *storage.Storage
	tokenAuth *jwtauth.JWTAuth
}

func New(cfg *config.Config, storage *storage.Storage) *Server {
	tokenAuth := jwtauth.New("HS256", []byte(cfg.App.SecretKey), nil)
	return &Server{cfg, storage, tokenAuth}
}

func (s *Server) Start() {
	router := s.createRouter()

	srv := &http.Server{
		Addr:    s.cfg.App.Address + ":" + s.cfg.App.Port,
		Handler: router,
	}

	serverErr := make(chan error, 1)

	go func() {
		fmt.Printf("Server starting on %s\n", srv.Addr)
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

	authHandler := authhandler.New(s.cfg, s.storage, slog.Default())
	auth := chi.NewRouter()
	auth.Post("/sign_up", authHandler.SignUp)
	auth.Post("/login", authHandler.SignIn)

	protected := chi.NewRouter()
	protected.Use(jwtauth.Verifier(s.tokenAuth))
	protected.Use(jwtauth.Authenticator(s.tokenAuth))

	wishlisthandler := wishlist.New(s.cfg, s.storage, slog.Default())
	protected.Get("/wishlist/{user_id}", wishlisthandler.GetWishlist)

	router.Mount("/auth", auth)
	router.Mount("/api", protected)

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
