package fleet_collector

import (
	"github.com/david-vtuk/prometheus-rancher-exporter/query/fleet"
	"github.com/ebauman/prometheus-rancher-exporter/query/rancher"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"time"
)

type metrics struct {
	fleetClusterGroupCount prometheus.Gauge
}

func new() metrics {
	m := metrics{
		fleetClusterGroupCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "fleet_cluster_groups",
			Help: "number of fleet cluster groups",
		}),
	}

	prometheus.MustRegister(m.fleetClusterGroupCount)

	m.fleetClusterGroupCount.Set(0)

	return m
}

func Collect(client rancher.Client) {

	log.Info("updating metrics")
	ticker := time.NewTicker(3 * time.Second)

	for range ticker.C {

		clustergroups, err := client.

	}

}