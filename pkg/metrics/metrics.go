package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	messagesProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "notification_messages_processed_total",
			Help: "Total number of messages processed from Kafka",
		},
	)
	messagesSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "notification_messages_success_total",
			Help: "Total number of successful Telegram notifications",
		},
	)
	messagesFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "notification_messages_failed_total",
			Help: "Total number of failed Telegram notifications",
		},
	)
	processingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "notification_processing_duration_seconds",
			Help:    "Time taken to process a message",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(messagesProcessed, messagesSuccess, messagesFailed, processingDuration)
}

func NewHandler() http.Handler {
	return promhttp.Handler()
}

func IncProcessed() {
	messagesProcessed.Inc()
}

func IncSuccess() {
	messagesSuccess.Inc()
}

func IncFailed() {
	messagesFailed.Inc()
}

func ObserveDuration(d time.Duration) {
	processingDuration.Observe(d.Seconds())
}