package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/client"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/config"
	api "github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/gen"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	runMode        *string
	cfg            *config.Config
	tokenAuth      *jwtauth.JWTAuth
	log            *slog.Logger
	authClient     *client.AuthClient
	wishClient     *client.WishlisterClient
	requestCounter prometheus.Counter
	responseTime   prometheus.Histogram
	errorCounter   prometheus.Counter
	activeRequests prometheus.Gauge
}

func New(runMode *string, cfg *config.Config, log *slog.Logger, authClient *client.AuthClient, wishClient *client.WishlisterClient) *Server {
	tokenAuth := jwtauth.New("HS256", []byte(cfg.App.SecretKey), nil)
	server := &Server{runMode, cfg, tokenAuth, log, authClient, wishClient, nil, nil, nil, nil}
	server.initMetrics()
	return server
}

func (s *Server) Start() {
	router := s.createRouter()

	srv := &http.Server{
		Addr:    s.cfg.App.Address + ":" + s.cfg.App.Port,
		Handler: router,
	}

	serverErr := make(chan error, 1)

	go func() {
		if *s.runMode != "production" {
			s.log.Info("Gateway server started",
				"address", s.cfg.App.Address,
				"port", s.cfg.App.Port,
				"mode", *s.runMode,
				"metrics", fmt.Sprintf("http://%s:%s/metrics", s.cfg.App.Address, s.cfg.App.Port))
		} else {
			s.log.Info("Gateway server started",
				"address", s.cfg.App.Address,
				"port", s.cfg.App.Port,
				"mode", *s.runMode)
		}

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
		close(serverErr)
	}()

	quit := make(chan os.Signal, 1)

	select {
	case <-quit:
		s.log.Info("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			s.log.Error("Server forced to shutdown", "error", err)
		}
		s.log.Info("Server exited properly")
	case err := <-serverErr:
		s.log.Error("Server error", "error", err)
	}
}

func (s *Server) createRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(corsMiddleware)
	router.Use(s.metricsMiddleware)

	if *s.runMode != "production" {
		router.Handle("/metrics", promhttp.Handler())
		router.Mount("/docs", s.redocRoutes())
		s.log.Info("Redoc documentation available",
			"address", s.cfg.App.Address,
			"port", s.cfg.App.Port,
			"url", fmt.Sprintf("http://%s:%s/docs", s.cfg.App.Address, s.cfg.App.Port))
	}
	router.Mount("/api/v1", s.apiRoutes())

	return router
}

func (s *Server) apiRoutes() http.Handler {
	r := chi.NewRouter()
	apiImpl := handlers.NewAPIImplementation(s.authClient, s.wishClient, s.log, s.tokenAuth)

	r.Group(func(r chi.Router) {
		r.Post("/auth/login", apiImpl.PostAuthLogin)
		r.Post("/auth/signup", apiImpl.PostAuthSignup)
		r.Get("/wishes/shared", apiImpl.GetWishesShared)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(s.tokenAuth))
		r.Use(jwtauth.Authenticator(s.tokenAuth))

		apiHandler := api.Handler(apiImpl)
		r.Mount("/", apiHandler)
	})

	return r
}
