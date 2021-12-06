package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	compbasemetrics "k8s.io/component-base/metrics"
)

var (
	defaultRegistry compbasemetrics.KubeRegistry
	// MustRegister registers registerable metrics but uses the defaultRegistry, panic upon the first registration that causes an error
	MustRegister func(...compbasemetrics.Registerable)
	// Register registers a collectable metric but uses the defaultRegistry
	Register func(compbasemetrics.Registerable) error

	RawMustRegister func(...prometheus.Collector)
)

func init() {
	defaultRegistry = compbasemetrics.NewKubeRegistry()
	MustRegister = defaultRegistry.MustRegister
	Register = defaultRegistry.Register
	RawMustRegister = defaultRegistry.RawMustRegister
}

func Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(prometheus.NewRegistry(), promhttp.HandlerFor(defaultRegistry, promhttp.HandlerOpts{}))
}
