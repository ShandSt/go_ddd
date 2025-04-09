package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	ProductOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_operations_total",
			Help: "Total number of product operations",
		},
		[]string{"operation", "status"},
	)

	ProductOperationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "product_operation_duration_seconds",
			Help:    "Duration of product operations in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"operation"},
	)

	MongoDBOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mongodb_operations_total",
			Help: "Total number of MongoDB operations",
		},
		[]string{"operation"},
	)

	MongoDBOperationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "mongodb_operation_duration_seconds",
			Help:    "Duration of MongoDB operations in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"operation"},
	)
)

func init() {
	prometheus.MustRegister(HTTPRequestsTotal)
	prometheus.MustRegister(HTTPRequestDuration)
	prometheus.MustRegister(ProductOperationsTotal)
	prometheus.MustRegister(ProductOperationDuration)
	prometheus.MustRegister(MongoDBOperationsTotal)
	prometheus.MustRegister(MongoDBOperationDuration)
}

func Handler() http.Handler {
	return promhttp.Handler()
}
