// Package prom implements....
package prom

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func P() {
	rege := prometheus.NewRegistry()
	des := make(chan<- *prometheus.Desc)
	rege.Describe(des)
	x := prometheus.CounterOpts{
		Name:        "counter",
		ConstLabels: prometheus.Labels{},
		Namespace:   "/",
		Subsystem:   "windows",
	}

	rege.Register(prometheus.NewCounter(x))
	gat, err := rege.Gather()
	if err != nil {
		panic(err)
	}
	for _, r := range gat {
		if r != nil {
			for _, m := range r.Metric {
				log.Printf(`
		guage %.6d,
       counter: %.6d,
      histogram: %v
		`, m.Gauge.Value, m.Counter.Value, m.Histogram)
			}
		}
	}
	http.Handle("/metrics", promhttp.HandlerFor(rege, promhttp.HandlerOpts{}))
	http.ListenAndServe(":2112", nil)
}
