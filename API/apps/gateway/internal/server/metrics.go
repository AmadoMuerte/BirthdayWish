package server

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
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
