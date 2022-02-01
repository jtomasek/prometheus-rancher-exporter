package fleet_collector

import (
	"github.com/david-vtuk/prometheus-rancher-exporter/query/fleet"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
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

func Collect(client fleet.Client) {

	m := new()

	fleetTicker := time.NewTicker(2 * time.Second)

	for range fleetTicker.C {
		log.Info("updating fleet metrics")
		clusterGroups, err := client.GetNumberOfClusterGroups()

		if err != nil {
			log.Errorf("error retrieving number of fleet cluster groups: %v", err)
		}

		m.fleetClusterGroupCount.Set(float64(clusterGroups))
	}

}
