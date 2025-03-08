package monitoring

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// CreateMetricsServer creates HTTP server for exporting prometheus metrics
func CreateMetricsServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return server
}
