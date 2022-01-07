package collector

import (
	"github.com/ebauman/prometheus-rancher-exporter/query/rancher"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"time"
)

type metrics struct {
	rancherMajorVersion prometheus.Gauge
	rancherMinorVersion prometheus.Gauge
	rancherPatchVersion prometheus.Gauge
	managedClusterCount prometheus.Gauge
}

func new() metrics {
	m := metrics{
		rancherMajorVersion: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rancher_major_version",
			Help: "major version for rancher installation, where version is semantic and formatted major.minor.patch",
		}),
		rancherMinorVersion: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rancher_minor_version",
			Help: "minor version for rancher installation, where version is semantic and formatted major.minor.patch",
		}),
		rancherPatchVersion: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rancher_patch_version",
			Help: "patch version for rancher installation, where version is semantic and formatted major.minor.patch",
		}),
		managedClusterCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rancher_managed_clusters",
			Help: "number of clusters this Rancher instance is currently managing",
		}),
	}

	prometheus.MustRegister(m.rancherMajorVersion)
	prometheus.MustRegister(m.rancherMinorVersion)
	prometheus.MustRegister(m.rancherPatchVersion)
	prometheus.MustRegister(m.managedClusterCount)

	m.rancherMajorVersion.Set(0)
	m.rancherMinorVersion.Set(0)
	m.rancherPatchVersion.Set(0)
	m.managedClusterCount.Set(0)

	return m
}

func Collect(client rancher.Client) {
	m := new()
	ticker := time.NewTicker(30 * time.Second)

	for range ticker.C {
		log.Info("updating metrics")

		vers, err := client.GetRancherVersion()

		if err != nil {
			log.Errorf("error retrieving rancher version: %v", err)
		}

		numberOfClusters, err := client.GetNumberOfManagedClusters()

		m.rancherMajorVersion.Set(float64(vers["major"]))
		m.rancherMinorVersion.Set(float64(vers["minor"]))
		m.rancherPatchVersion.Set(float64(vers["patch"]))

		m.managedClusterCount.Set(float64(numberOfClusters))
	}

}
