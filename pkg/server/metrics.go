package server

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "stata"
)

var (
	metricStatusCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_status_codes_count",
			Help:      "Number of HTTP requests by status code.",
		},
		[]string{"code"},
	)

	metricsEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "bot_events_count",
			Help:      "Number of events by type.",
		},
		[]string{"bot_token", "event_type"},
	)
)

func init() {
	prometheus.MustRegister(metricStatusCodes)
	prometheus.MustRegister(metricsEvents)
}
