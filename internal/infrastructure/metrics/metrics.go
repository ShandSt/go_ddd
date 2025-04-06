package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestsTotal tracks the total number of HTTP requests
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestDuration tracks the duration of HTTP requests
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// ProductOperationsTotal tracks the total number of product operations
	ProductOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_operations_total",
			Help: "Total number of product operations",
		},
		[]string{"operation", "status"},
	)

	// ProductOperationDuration tracks the duration of product operations
	ProductOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "product_operation_duration_seconds",
			Help:    "Duration of product operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// MongoDBOperationsTotal tracks the total number of MongoDB operations
	MongoDBOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mongodb_operations_total",
			Help: "Total number of MongoDB operations",
		},
		[]string{"operation", "status"},
	)

	// MongoDBOperationDuration tracks the duration of MongoDB operations
	MongoDBOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "mongodb_operation_duration_seconds",
			Help:    "Duration of MongoDB operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)
