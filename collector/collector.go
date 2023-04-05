package collector

import (
	"github.com/david-vtuk/prometheus-rancher-exporter/query/rancher"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"time"
)

type metrics struct {
	installedRancherVersion prometheus.GaugeVec
	latestRancherVersion    prometheus.GaugeVec

	managedClusterCount     prometheus.Gauge
	managedK3sClusterCount  prometheus.Gauge
	managedRKEClusterCount  prometheus.Gauge
	managedRKE2ClusterCount prometheus.Gauge
	managedEKSClusterCount  prometheus.Gauge
	managedAKSClusterCount  prometheus.Gauge
	managedGKEClusterCount  prometheus.Gauge
	managedNodeCount        prometheus.Gauge

	// Cluster level metrics
	clusterConditionConnected    prometheus.GaugeVec
	clusterConditionNotConnected prometheus.GaugeVec

	// Downstream cluster metrics
	downstreamClusterVersion prometheus.GaugeVec

	// User related
	tokenCount prometheus.Gauge
	userCount  prometheus.Gauge

	// Project related
	projectCount       prometheus.Gauge
	projectLabels      prometheus.GaugeVec
	projectAnnotations prometheus.GaugeVec
	projectResources   prometheus.GaugeVec
}

func new() metrics {
	m := metrics{

		installedRancherVersion: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "installed_rancher_version",
			Help: "version of the installed Rancher instance",
		}, []string{"version"},
		),
		latestRancherVersion: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "latest_rancher_version",
			Help: "version of the most recent Rancher release",
		}, []string{"version"},
		),
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
		managedNodeCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rancher_managed_nodes",
			Help: "number of managed nodes this Rancher instance is currently managing",
		}),
		clusterConditionConnected: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cluster_connected",
			Help: "identify if a downstream cluster is connected to Rancher",
		}, []string{"Name"},
		),
		clusterConditionNotConnected: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cluster_not_connected",
			Help: "identify if a downstream cluster is not connected to Rancher",
		}, []string{"Name"},
		),
		downstreamClusterVersion: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cluster_k8s_version",
			Help: "version of K8s running in the downstream cluster",
		}, []string{"Name", "Version"},
		),
		tokenCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rancher_tokens",
			Help: "number of tokens issued by Rancher",
		}),
		userCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rancher_users",
			Help: "number of users in this Rancher instance",
		}),
		projectCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rancher_projects",
			Help: "number of Projects globally",
		}),
		projectLabels: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "rancher_project_labels",
			Help: "labels associated with Rancher Projects",
		}, []string{"cluster_name", "project_id", "project_display_name", "project_label_key", "project_label_value"},
		),
		projectAnnotations: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "rancher_project_annotations",
			Help: "annotations associated with Rancher Projects",
		}, []string{"cluster_name", "project_id", "project_display_name", "project_annotation_key", "project_annotation_value"},
		),
		projectResources: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "rancher_project_resourcequota",
			Help: "default namespace resource quota set the for project",
		}, []string{"cluster_name", "project_id", "project_display_name", "project_resource_key", "project_resource_type"},
		),
	}

	prometheus.MustRegister(m.installedRancherVersion)
	prometheus.MustRegister(m.latestRancherVersion)
	prometheus.MustRegister(m.managedClusterCount)
	prometheus.MustRegister(m.managedRKEClusterCount)
	prometheus.MustRegister(m.managedRKE2ClusterCount)
	prometheus.MustRegister(m.managedK3sClusterCount)
	prometheus.MustRegister(m.managedEKSClusterCount)
	prometheus.MustRegister(m.managedAKSClusterCount)
	prometheus.MustRegister(m.managedGKEClusterCount)
	prometheus.MustRegister(m.managedNodeCount)
	prometheus.MustRegister(m.clusterConditionConnected)
	prometheus.MustRegister(m.clusterConditionNotConnected)

	prometheus.MustRegister(m.downstreamClusterVersion)

	prometheus.MustRegister(m.tokenCount)
	prometheus.MustRegister(m.userCount)

	prometheus.MustRegister(m.projectCount)
	prometheus.MustRegister(m.projectLabels)
	prometheus.MustRegister(m.projectAnnotations)
	prometheus.MustRegister(m.projectResources)

	m.managedClusterCount.Set(0)
	m.managedRKEClusterCount.Set(0)
	m.managedRKE2ClusterCount.Set(0)
	m.managedK3sClusterCount.Set(0)
	m.managedEKSClusterCount.Set(0)
	m.managedAKSClusterCount.Set(0)
	m.managedGKEClusterCount.Set(0)
	m.managedNodeCount.Set(0)

	m.clusterConditionConnected.Reset()
	m.clusterConditionNotConnected.Reset()

	m.downstreamClusterVersion.Reset()

	m.tokenCount.Set(0)
	m.userCount.Set(0)

	m.projectCount.Set(0)
	m.projectLabels.Reset()

	return m
}

func Collect(client rancher.Client) {
	m := new()

	// GitHub API request limits necessitate polling at a different interval

	go func() {
		ticker := time.NewTicker(1 * time.Minute)

		for range ticker.C {

			m.latestRancherVersion.Reset()

			latestVers, err := client.GetLatestRancherVersion()

			if err != nil {
				log.Errorf("error retrieving latest Rancher version: %v", err)
			}

			m.latestRancherVersion.WithLabelValues(latestVers).Set(1)
		}
	}()

	ticker := time.NewTicker(3 * time.Second)

	for range ticker.C {
		log.Info("updating rancher metrics")

		resetGaugeMetrics(m)

		installedVersion, err := client.GetInstalledRancherVersion()
		if err != nil {
			return
		}

		state, err := client.GetClusterConnectedState()
		if err != nil {
			return
		}

		numberOfClusters, err := client.GetNumberOfManagedClusters()

		if err != nil {
			log.Errorf("error retrieving number of managed clusters: %v", err)
		}

		distributions, err := client.GetK8sDistributions()

		if err != nil {
			log.Errorf("error retrieving number of managed clusters: %v", err)
		}

		numberOfNodes, err := client.GetNumberOfManagedNodes()

		if err != nil {
			log.Errorf("error retrieving number of managed nodes: %v", err)
		}

		downstreamClusterVersions, err := client.GetDownstreamClusterVersions()

		if err != nil {
			log.Errorf("error retrieving downstream k8s cluster versions: %v", err)
		}

		m.installedRancherVersion.WithLabelValues(installedVersion).Set(1)

		for _, value := range downstreamClusterVersions {

			m.downstreamClusterVersion.WithLabelValues(value.Name, value.Version).Set(1)

		}

		m.managedClusterCount.Set(float64(numberOfClusters))
		m.managedNodeCount.Set(float64(numberOfNodes))

		m.managedRKEClusterCount.Set(float64(distributions["rke"]))
		m.managedRKE2ClusterCount.Set(float64(distributions["rke2"]))
		m.managedK3sClusterCount.Set(float64(distributions["k3s"]))
		m.managedEKSClusterCount.Set(float64(distributions["eks"]))
		m.managedAKSClusterCount.Set(float64(distributions["aks"]))
		m.managedGKEClusterCount.Set(float64(distributions["gke"]))

		for key, value := range state {
			if value == true {
				m.clusterConditionConnected.WithLabelValues(key).Set(1)
				m.clusterConditionNotConnected.WithLabelValues(key).Set(0)
			} else {
				m.clusterConditionNotConnected.WithLabelValues(key).Set(1)
				m.clusterConditionConnected.WithLabelValues(key).Set(0)
			}
		}

		tokens, err := client.GetNumberOfTokens()
		if err != nil {
			log.Errorf("error retrieving number of tokens: %v", err)
		}

		m.tokenCount.Set(float64(tokens))

		users, err := client.GetNumberOfUsers()
		if err != nil {
			log.Errorf("error retrieving number of users: %v", err)
		}

		m.userCount.Set(float64(users))

		projects, err := client.GetNumberofProjects()
		if err != nil {
			log.Errorf("error retrieving number of projects: %v", err)
		}

		m.projectCount.Set(float64(projects))

		projectLabels, err := client.GetProjectLabels()
		if err != nil {
			log.Errorf("error retrieving project labels: %v", err)
		}

		for _, value := range projectLabels {

			m.projectLabels.WithLabelValues(value.ProjectClusterName, value.Projectid, value.ProjectDisplayName, value.LabelKey, value.LabelValue).Set(1)

		}

		projectAnnotations, err := client.GetProjectAnnotations()
		if err != nil {
			log.Errorf("error retrieving project annotations: %v", err)
		}

		for _, value := range projectAnnotations {
			m.projectAnnotations.WithLabelValues(value.ProjectClusterName, value.Projectid, value.ProjectDisplayName, value.AnnotationKey, value.AnnotationValue).Set(1)
		}

		projectResources, err := client.GetProjectResourceQuota()
		for _, value := range projectResources {
			m.projectResources.WithLabelValues(value.ProjectClusterName, value.Projectid, value.ProjectDisplayName, value.ResourceKey, value.ResourceType).Set(value.ResourceValue)
		}
	}

}

// Reset GaugeVecs on each tick - facilitate state transition
func resetGaugeMetrics(m metrics) {

	m.downstreamClusterVersion.Reset()
	m.clusterConditionNotConnected.Reset()
	m.clusterConditionConnected.Reset()
	m.installedRancherVersion.Reset()
}
