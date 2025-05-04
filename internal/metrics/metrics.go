package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	// TCP Gateway metrics
	TCPConnectionsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tcp_connections_total",
		Help: "The total number of TCP connections received",
	})

	TCPConnectionErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tcp_connection_errors_total",
		Help: "The total number of TCP connection errors",
	})

	TCPBytesReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tcp_bytes_received_total",
		Help: "The total number of bytes received via TCP",
	})

	// Parser metrics
	ParserProcessedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "parser_processed_total",
		Help: "The total number of messages processed by the parser",
	})

	ParserErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "parser_errors_total",
		Help: "The total number of parser errors",
	})

	ParserLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "parser_latency_seconds",
		Help:    "The time taken to parse messages",
		Buckets: prometheus.DefBuckets,
	})

	// MQTT Publisher metrics
	MQTTPublishedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_published_total",
		Help: "The total number of messages published to MQTT",
	})

	MQTTPublishErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_publish_errors_total",
		Help: "The total number of MQTT publish errors",
	})

	MQTTPublishLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "mqtt_publish_latency_seconds",
		Help:    "The time taken to publish messages to MQTT",
		Buckets: prometheus.DefBuckets,
	})

	// MQTT Subscriber metrics
	MQTTSubscriberReceivedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_subscriber_received_total",
		Help: "The total number of messages received by MQTT subscriber",
	})

	// Storage metrics
	StorageOperationsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "storage_operations_total",
		Help: "The total number of storage operations",
	})

	StorageErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "storage_errors_total",
		Help: "The total number of storage errors",
	})

	StorageLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "storage_latency_seconds",
		Help:    "The time taken to store messages",
		Buckets: prometheus.DefBuckets,
	})

	// Processor metrics
	ProcessorProcessedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "processor_processed_total",
		Help: "The total number of messages processed",
	})

	ProcessorErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "processor_errors_total",
		Help: "The total number of processor errors",
	})

	ProcessorLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "processor_latency_seconds",
		Help:    "The time taken to process messages",
		Buckets: prometheus.DefBuckets,
	})

	// Security metrics
	FailedLoginAttemptsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "failed_login_attempts_total",
		Help: "The total number of failed login attempts",
	})

	UnauthorizedAccessAttemptsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "unauthorized_access_attempts_total",
		Help: "The total number of unauthorized access attempts",
	})
)

// RecordParserLatency records the time taken to parse a message
func RecordParserLatency(start time.Time) {
	ParserLatency.Observe(time.Since(start).Seconds())
}

// RecordMQTTPublishLatency records the time taken to publish a message to MQTT
func RecordMQTTPublishLatency(start time.Time) {
	MQTTPublishLatency.Observe(time.Since(start).Seconds())
}

// RecordStorageLatency records the time taken to store a message
func RecordStorageLatency(start time.Time) {
	StorageLatency.Observe(time.Since(start).Seconds())
}

// RecordProcessorLatency records the time taken to process a message
func RecordProcessorLatency(start time.Time) {
	ProcessorLatency.Observe(time.Since(start).Seconds())
}

// StartMetricsServer starts an HTTP server to expose Prometheus metrics
func StartMetricsServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
		if err != nil {
			panic(err)
		}
	}()
}
