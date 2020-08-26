package monitor

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMonitor struct {
	NoEnvTagWatchdog *prometheus.GaugeVec
	NoEnvTag         *prometheus.GaugeVec
}

func NewPrometheusMonitor() *PrometheusMonitor {
	NoEnvTagWatchdog := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "notagwatchdog",
			Help: "watchdog for no tag checking program.",
		},
		[]string{"name"},
	)
	NoEnvTag := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "notag",
			Help: "No environment tag on ecs instance.",
		},
		[]string{"id", "vpc", "name"},
	)

	prometheus.MustRegister(NoEnvTagWatchdog)
	prometheus.MustRegister(NoEnvTag)

	return &PrometheusMonitor{
		NoEnvTagWatchdog: NoEnvTagWatchdog,
		NoEnvTag:         NoEnvTag,
	}
}

func PrometheusBoot() error {
	http.Handle("/metrics", promhttp.Handler())
	return func() error {
		err := http.ListenAndServe("0.0.0.0:8085", nil)
		if err != nil {
			return err
		}
		return nil
	}()
}
