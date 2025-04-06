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
	"github.com/AmadoMuerte/BirthdayWish/API/internal/storage"
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

	authHandler := authhandler.New(s.cfg, s.storage, slog.Default())
	// noteHandler := notehandler.New(s.cfg, s.db, log)
	// userHandler := userhandler.New(s.cfg, s.db, log)

	auth := chi.NewRouter()
	// auth.Post("/sign-in", authHandler.SignIn)
	auth.Post("/sign_up", authHandler.SignUp)
	// auth.Post("/refresh", authHandler.Refresh)

	// user := chi.NewRouter()
	// user.Use(middlewares.AuthMiddleware)
	// user.Get("/{id}", userHandler.Get)
	// user.Put("/", userHandler.Update)
	// user.Delete("/{id}", userHandler.Delete)

	// note := chi.NewRouter()
	// note.Use(middlewares.AuthMiddleware)
	// note.Get("/", noteHandler.Get)
	// note.Post("/", noteHandler.Post)

	// router.Mount("/note", note)
	router.Mount("/auth", auth)
	// router.Mount("/users", user)

	return router
}
