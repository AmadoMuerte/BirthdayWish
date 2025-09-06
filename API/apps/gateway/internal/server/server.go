package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/config"
	api "github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/gen"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/handlers"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mvrilo/go-redoc"
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
		fmt.Printf("Gateway starting on %s\n", srv.Addr)
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

	router.Mount("/auth", s.authRoutes())
	router.Mount("/api", s.apiRoutes())

	if s.cfg.App.Mode == "dev" {
		router.Mount("/docs", s.redocRoutes())
	}

	return router
}

func (s *Server) redocRoutes() http.Handler {
	r := chi.NewRouter()

	doc := redoc.Redoc{
		Title:       "BirthdayWish API",
		Description: "Gateway API for BirthdayWish",
		SpecFile:    "internal/api/openapi.yaml",
		SpecPath:    "/docs/openapi.yaml",
		DocsPath:    "/docs",
	}

	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/api/openapi.yaml")
	})

	r.Get("/", doc.Handler())

	return r
}

func (s *Server) authRoutes() http.Handler {
	r := chi.NewRouter()
	// authHandler := authhandler.New(s.cfg, s.storage, slog.Default())
	// wishhandler := wishlist.New(s.cfg, s.storage, s.RedisClient, slog.Default())

	// r.Post("/sign_up", authHandler.SignUp)
	// r.Post("/login", authHandler.SignIn)
	// r.Get("/get_wishlist", wishhandler.GetShareList)
	return r
}

func (s *Server) apiRoutes() http.Handler {
	r := chi.NewRouter()

	apiImpl := handlers.NewAPIImplementation()
	apiHandler := api.Handler(apiImpl)

	r.Mount("/v1", apiHandler)

	protected := chi.NewRouter()
	protected.Use(jwtauth.Verifier(s.tokenAuth))
	protected.Use(jwtauth.Authenticator(s.tokenAuth))

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
