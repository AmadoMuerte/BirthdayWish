package server

import (
	"context"
	"fmt"
	"log/slog"
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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	cfg            *config.Config
	storage        *storage.Storage
	tokenAuth      *jwtauth.JWTAuth
	log            *slog.Logger
	requestCounter prometheus.Counter
	responseTime   prometheus.Histogram
	errorCounter   prometheus.Counter
	activeRequests prometheus.Gauge
}

func New(cfg *config.Config, storage *storage.Storage, log *slog.Logger) *Server {
	tokenAuth := jwtauth.New("HS256", []byte(cfg.App.SecretKey), nil)
	server := &Server{cfg, storage, tokenAuth, log, nil, nil, nil, nil}
	server.initMetrics()
	return server
}

func (s *Server) initMetrics() {
	s.requestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
		ConstLabels: prometheus.Labels{
			"service": "gateway",
			"version": "1.0",
		},
	})

	s.responseTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
		ConstLabels: prometheus.Labels{
			"service": "gateway",
		},
	})

	s.errorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "Total number of HTTP errors",
		ConstLabels: prometheus.Labels{
			"service": "gateway",
		},
	})

	s.activeRequests = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_active",
		Help: "Number of active HTTP requests",
		ConstLabels: prometheus.Labels{
			"service": "gateway",
		},
	})
}

func (s *Server) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		s.activeRequests.Inc()
		defer s.activeRequests.Dec()
		s.requestCounter.Inc()

		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		s.responseTime.Observe(duration)

		if rw.statusCode >= 400 {
			s.errorCounter.Inc()
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (s *Server) Start() {
	router := s.createRouter()

	srv := &http.Server{
		Addr:    s.cfg.App.Address + ":" + s.cfg.App.Port,
		Handler: router,
	}

	serverErr := make(chan error, 1)

	go func() {
		s.log.Info("Gateway server started",
			"address", s.cfg.App.Address,
			"port", s.cfg.App.Port,
			"metrics", fmt.Sprintf("http://%s:%s/metrics", s.cfg.App.Address, s.cfg.App.Port))

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

	router.Handle("/metrics", promhttp.Handler())

	router.Mount("/auth", s.authRoutes())
	router.Mount("/api", s.apiRoutes())

	if s.cfg.App.Mode == "dev" {
		router.Mount("/docs", s.redocRoutes())
		s.log.Info("Redoc documentation available",
			"address", s.cfg.App.Address,
			"port", s.cfg.App.Port,
			"url", fmt.Sprintf("http://%s:%s/docs", s.cfg.App.Address, s.cfg.App.Port))
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
