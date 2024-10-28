package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
)

var (
	HttpRequestCountWithPath = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total_with_path",
			Help: "Number of HTTP requests by path.",
		},
		[]string{"url"},
	)

	// PROMQL => rate(http_request_duration_seconds_sum{}[5m]) / rate(http_request_duration_seconds_count{}[5m])
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Response time of HTTP request.",
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(HttpRequestCountWithPath)
	prometheus.MustRegister(HttpRequestDuration)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	r.HandleFunc("/register", registerHandler).Methods("POST")
	r.HandleFunc("/deregister", deregisterHandler).Methods("DELETE")
	r.HandleFunc("/discover", discoverHandler).Methods("GET")

	http.ListenAndServe(":8000", r)
}
