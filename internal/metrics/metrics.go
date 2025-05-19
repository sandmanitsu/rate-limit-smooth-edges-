package metric

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var templateMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "api",
	Subsystem:  "request",
	Name:       "code",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"code"})

func ObserveCodeStatus(code int, d time.Duration) {
	templateMetrics.WithLabelValues(strconv.Itoa(code)).Observe(d.Seconds())
}

func Listen(host string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(host, mux)
}
