package collector

import (
	"github.com/ebauman/prometheus-rancher-exporter/query/rancher"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"time"
)

type metrics struct {
	rancherMajorVersion     prometheus.Gauge
	rancherMinorVersion     prometheus.Gauge
	rancherPatchVersion     prometheus.Gauge
	managedClusterCount     prometheus.Gauge
	managedK3sClusterCount  prometheus.Gauge
	managedRKEClusterCount  prometheus.Gauge
	managedRKE2ClusterCount prometheus.Gauge
	managedEKSClusterCount  prometheus.Gauge
	managedAKSClusterCount  prometheus.Gauge
	managedGKEClusterCount  prometheus.Gauge
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
		managedRKEClusterCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rancher_managed_rke_clusters",
			Help: "number of RKE clusters this Rancher instance is currently managing",
		}),
		managedRKE2ClusterCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rancher_managed_rke2_clusters",
			Help: "number of RKE2 clusters this Rancher instance is currently managing",
		}),
		managedK3sClusterCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rancher_managed_k3s_clusters",
			Help: "number of K3s clusters this Rancher instance is currently managing",
		}),
		managedEKSClusterCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rancher_managed_eks_clusters",
			Help: "number of RKE clusters this Rancher instance is currently managing",
		}),
		managedAKSClusterCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rancher_managed_aks_clusters",
			Help: "number of RKE clusters this Rancher instance is currently managing",
		}),
		managedGKEClusterCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rancher_managed_gke_clusters",
			Help: "number of RKE clusters this Rancher instance is currently managing",
		}),
	}

	prometheus.MustRegister(m.rancherMajorVersion)
	prometheus.MustRegister(m.rancherMinorVersion)
	prometheus.MustRegister(m.rancherPatchVersion)
	prometheus.MustRegister(m.managedClusterCount)
	prometheus.MustRegister(m.managedRKEClusterCount)
	prometheus.MustRegister(m.managedRKE2ClusterCount)
	prometheus.MustRegister(m.managedK3sClusterCount)
	prometheus.MustRegister(m.managedEKSClusterCount)
	prometheus.MustRegister(m.managedAKSClusterCount)
	prometheus.MustRegister(m.managedGKEClusterCount)

	m.rancherMajorVersion.Set(0)
	m.rancherMinorVersion.Set(0)
	m.rancherPatchVersion.Set(0)
	m.managedClusterCount.Set(0)
	m.managedRKEClusterCount.Set(0)
	m.managedRKE2ClusterCount.Set(0)
	m.managedK3sClusterCount.Set(0)
	m.managedEKSClusterCount.Set(0)
	m.managedAKSClusterCount.Set(0)
	m.managedGKEClusterCount.Set(0)

	return m
}

func Collect(client rancher.Client) {
	m := new()
	ticker := time.NewTicker(3 * time.Second)

	for range ticker.C {
		log.Info("updating metrics")

		vers, err := client.GetRancherVersion()

		if err != nil {
			log.Errorf("error retrieving rancher version: %v", err)
		}

		numberOfClusters, err := client.GetNumberOfManagedClusters()

		if err != nil {
			log.Errorf("error retrieving number of managed clusters: %v", err)
		}

		distributions, err := client.GetK8sDistributions()

		if err != nil {
			log.Errorf("error retrieving number of managed clusters: %v", err)
		}

		m.rancherMajorVersion.Set(float64(vers["major"]))
		m.rancherMinorVersion.Set(float64(vers["minor"]))
		m.rancherPatchVersion.Set(float64(vers["patch"]))

		m.managedClusterCount.Set(float64(numberOfClusters))

		m.managedRKEClusterCount.Set(float64(distributions["rke"]))
		m.managedRKE2ClusterCount.Set(float64(distributions["rke2"]))
		m.managedK3sClusterCount.Set(float64(distributions["k3s"]))
		m.managedEKSClusterCount.Set(float64(distributions["eks"]))
		m.managedAKSClusterCount.Set(float64(distributions["aks"]))
		m.managedGKEClusterCount.Set(float64(distributions["gke"]))
	}

}
