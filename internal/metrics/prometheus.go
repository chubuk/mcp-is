package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	DBQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
		},
		[]string{"operation"},
	)

	MCPToolCalls = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mcp_tool_calls_total",
			Help: "Total MCP tool calls",
		},
		[]string{"tool", "status"},
	)

	MCPToolDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "mcp_tool_duration_seconds",
			Help:    "MCP tool execution duration",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 5},
		},
		[]string{"tool"},
	)

	ActiveConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)

	StartupTime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_startup_time_seconds",
			Help: "Application startup time in seconds",
		},
	)
)

func init() {
	prometheus.MustRegister(
		HTTPRequests,
		HTTPRequestDuration,
		DBQueryDuration,
		MCPToolCalls,
		MCPToolDuration,
		ActiveConnections,
		StartupTime,
	)
}

// Handler returns the /metrics endpoint handler
func Handler() http.Handler {
	return promhttp.Handler()
}

// Timer is a helper for measuring duration
type Timer struct {
	start    time.Time
	duration func(float64)
}

// NewTimer starts a timer
func NewTimer(observer func(float64)) *Timer {
	return &Timer{start: time.Now(), duration: observer}
}

// Observe records the duration
func (t *Timer) Observe() {
	t.duration(time.Since(t.start).Seconds())
}
