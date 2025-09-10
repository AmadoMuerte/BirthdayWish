package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/service"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/storage"
	wishProto "github.com/AmadoMuerte/BirthdayWish/API/proto/wish"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	runMode        *string
	cfg            *config.Config
	storage        *storage.Storage
	wishService    *service.WishService
	grpcServer     *grpc.Server
	log            *slog.Logger
	requestCounter prometheus.Counter
	responseTime   prometheus.Histogram
	errorCounter   prometheus.Counter
}

func New(runMode *string, cfg *config.Config, storage *storage.Storage, wishService *service.WishService, log *slog.Logger) *Server {
	server := &Server{
		runMode:     runMode,
		cfg:         cfg,
		storage:     storage,
		wishService: wishService,
		log:         log,
	}
	server.initMetrics()
	return server
}

func (s *Server) initMetrics() {
	s.requestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_requests_total",
		Help: "Total number of gRPC requests",
		ConstLabels: prometheus.Labels{
			"service": "wish",
		},
	})

	s.responseTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "grpc_request_duration_seconds",
		Help:    "Duration of gRPC requests",
		Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2},
		ConstLabels: prometheus.Labels{
			"service": "wish",
		},
	})

	s.errorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_errors_total",
		Help: "Total number of gRPC errors",
		ConstLabels: prometheus.Labels{
			"service": "wish",
		},
	})
}

func (s *Server) Start() {
	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(s.metricsInterceptor),
	)

	wishProto.RegisterWishServiceServer(s.grpcServer, s.wishService)

	if *s.runMode != "production" {
		reflection.Register(s.grpcServer)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.cfg.App.Host, s.cfg.App.Port))
	if err != nil {
		s.log.Error("failed to listen", "error", err, "port", s.cfg.App.Port)
		os.Exit(1)
	}

	serverErr := make(chan error, 1)

	go func() {
		s.log.Info("Wish service started",
			"host", s.cfg.App.Host,
			"port", s.cfg.App.Port,
			"mode", *s.runMode)

		if err := s.grpcServer.Serve(lis); err != nil {
			serverErr <- err
		}
		close(serverErr)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		s.log.Info("Shutting down wish service...")
		s.Stop()
	case err := <-serverErr:
		s.log.Error("Server error", "error", err)
		s.Stop()
	}
}

func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	s.log.Info("Wish service stopped")
}

func (s *Server) metricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	s.requestCounter.Inc()

	resp, err := handler(ctx, req)

	duration := time.Since(start).Seconds()
	s.responseTime.Observe(duration)

	if err != nil {
		s.errorCounter.Inc()
	}

	return resp, err
}
