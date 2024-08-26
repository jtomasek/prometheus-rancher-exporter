package collector

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/david-vtuk/prometheus-rancher-exporter/query/rancher"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type rancherMetrics struct {
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

	// Cluster level rancherMetrics
	clusterConditionConnected      prometheus.GaugeVec
	clusterConditionNotConnected   prometheus.GaugeVec
	clusterConditionPending        prometheus.GaugeVec
	clusterConditionWaiting        prometheus.GaugeVec
	clusterConditionDiskPressure   prometheus.GaugeVec
	clusterConditionMemoryPressure prometheus.GaugeVec
	clusterConditionReady          prometheus.GaugeVec

	// Downstream cluster rancherMetrics
	downstreamClusterVersion prometheus.GaugeVec
	downstreamClusterStatus  prometheus.GaugeVec

	// User related
	tokenCount prometheus.Gauge
	userCount  prometheus.Gauge

	// Project related
	projectCount       prometheus.Gauge
	projectLabels      prometheus.GaugeVec
	projectAnnotations prometheus.GaugeVec
	projectResources   prometheus.GaugeVec

	// Extended rancherMetrics for Rancher CR's
	rancherCustomResources prometheus.GaugeVec

	// Additional Node Information
	managedNodeInfo prometheus.GaugeVec
}

type rancherBackupMetrics struct {
	backupCount     prometheus.Gauge
	restoreCount    prometheus.Gauge
	restoreSetCount prometheus.Gauge
	backup          prometheus.GaugeVec
	restore         prometheus.GaugeVec
}

func initRancherMetrics() rancherMetrics {
	m := rancherMetrics{

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
		rancherCustomResources: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "rancher_custom_resource_count",
			Help: "raw count of Rancher custom resources by name",
		}, []string{"resource_name"}),
		managedNodeInfo: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "rancher_managed_nodes_info",
			Help: "additional metadata about known downstream cluster nodes",
		}, []string{"name", "parent_cluster", "is_control_plane", "is_etcd", "is_worker", "architecture",
			"container_runtime_version", "kernel_version", "os", "os_image"},
		),
		downstreamClusterStatus: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "rancher_cluster_status",
			Help: "Status of the downstream cluster",
		}, []string{"Name", "DisplayName", "Version", "RancherServerURL"},
		),
		clusterConditionPending: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "rancher_cluster_condition_pending",
			Help: "Identify if a downstream cluster is in Pending state",
		}, []string{"Name"},
		),
		clusterConditionWaiting: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "rancher_cluster_condition_waiting",
			Help: "Identify if a downstream cluster is in Waiting state",
		}, []string{"Name"},
		),
		clusterConditionDiskPressure: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "rancher_cluster_condition_disk_pressure",
			Help: "Identify if a downstream cluster is running out of disk",
		}, []string{"Name"},
		),
		clusterConditionMemoryPressure: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "rancher_cluster_condition_memory_pressure",
			Help: "Identify if a downstream cluster is running out of memory",
		}, []string{"Name"},
		),
		clusterConditionReady: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "rancher_cluster_condition_ready",
			Help: "Identify if a downstream cluster is in Ready state",
		}, []string{"Name", "Reason", "Message"},
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
	prometheus.MustRegister(m.clusterConditionPending)
	prometheus.MustRegister(m.clusterConditionWaiting)
	prometheus.MustRegister(m.clusterConditionDiskPressure)
	prometheus.MustRegister(m.clusterConditionMemoryPressure)
	prometheus.MustRegister(m.clusterConditionReady)

	prometheus.MustRegister(m.downstreamClusterVersion)
	prometheus.MustRegister(m.downstreamClusterStatus)

	prometheus.MustRegister(m.tokenCount)
	prometheus.MustRegister(m.userCount)

	prometheus.MustRegister(m.projectCount)
	prometheus.MustRegister(m.projectLabels)
	prometheus.MustRegister(m.projectAnnotations)
	prometheus.MustRegister(m.projectResources)

	prometheus.MustRegister(m.rancherCustomResources)
	prometheus.MustRegister(m.managedNodeInfo)

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
	m.clusterConditionPending.Reset()
	m.clusterConditionWaiting.Reset()
	m.clusterConditionDiskPressure.Reset()
	m.clusterConditionMemoryPressure.Reset()
	m.clusterConditionReady.Reset()

	m.downstreamClusterVersion.Reset()
	m.downstreamClusterStatus.Reset()

	m.tokenCount.Set(0)
	m.userCount.Set(0)

	m.projectCount.Set(0)
	m.managedNodeInfo.Reset()

	return m
}

func initRancherBackupMetrics(reg *prometheus.Registry) rancherBackupMetrics {

	rBackupMetrics := rancherBackupMetrics{
		backupCount: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Name: "rancher_backups_count",
			Help: "number of rancher backups",
		}),
		restoreCount: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Name: "rancher_restore_count",
			Help: "number of rancher restore",
		}),
		backup: *promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Name: "rancher_backup",
			Help: "details regarding a specific backup operation",
		}, []string{"name", "resourceSetName", "retentionCount", "backupType", "status", "filename", "storageLocation", "nextSnapshot", "lastSnapshot"},
		),
		restore: *promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Name: "rancher_restore",
			Help: "details regarding a specific restore operation",
		}, []string{"name", "fileName", "prune", "storageLocation", "status", "restoreTime"},
		),
	}

	rBackupMetrics.backupCount.Set(0)
	rBackupMetrics.restoreCount.Set(0)
	rBackupMetrics.backup.Reset()

	return rBackupMetrics
}

func Collect(client rancher.Client, Timer_GetLatestRancherVersion int, Timer_ticker int) {

	baseMetrics := initRancherMetrics()

	// GitHub API request limits necessitate polling at a different interval

	go func() {
		ticker := time.NewTicker(time.Duration(Timer_GetLatestRancherVersion) * time.Minute)

		for ; ; <-ticker.C {

			baseMetrics.latestRancherVersion.Reset()

			latestVers, err := client.GetLatestRancherVersion()

			if err != nil {
				log.Errorf("error retrieving latest Rancher version: %v", err)
			}

			baseMetrics.latestRancherVersion.WithLabelValues(latestVers).Set(1)
		}
	}()

	ticker := time.NewTicker(time.Duration(Timer_ticker) * time.Second)

	for ; ; <-ticker.C {

		resetGaugeVecMetrics(baseMetrics)
		log.Info("Updating Metrics")

		go getInstalledRancherVersion(client, baseMetrics)
		go getClusterConnectedState(client, baseMetrics)
		go getNumberOfClusters(client, baseMetrics)
		go getDistributions(client, baseMetrics)
		go getNumberOfNodes(client, baseMetrics)
		go getDownstreamClusterVersions(client, baseMetrics)
		go getDownstreamClusterStates(client, baseMetrics)
		go getClusterConditions(client, baseMetrics)
		go getNumberOfTokens(client, baseMetrics)
		go getNumberOfUsers(client, baseMetrics)
		go getNumberOfProjects(client, baseMetrics)
		go getProjectLabels(client, baseMetrics)
		go getProjectAnnotations(client, baseMetrics)
		go getProjectResources(client, baseMetrics)
		go getRancherCustomResources(client, baseMetrics)
		go getNodeInfo(client, baseMetrics)

	}
}

func CollectBackupMetrics(client rancher.Client, Timer_ticker int, reg *prometheus.Registry) {

	backupMetrics := initRancherBackupMetrics(reg)

	ticker := time.NewTicker(time.Duration(Timer_ticker) * time.Second)

	for ; ; <-ticker.C {

		resetBackupGaugeVecMetrics(backupMetrics)
		go getNumberOfBackups(client, backupMetrics)
		go getNumberOfRestores(client, backupMetrics)
		go getBackups(client, backupMetrics)
		go getRestores(client, backupMetrics)
	}

}

func getInstalledRancherVersion(client rancher.Client, m rancherMetrics) {

	start := time.Now()

	installedVersion, err := client.GetInstalledRancherVersion()
	if err != nil {
		log.Errorf("error retrieving the installed Rancher version: %v", err)
	}
	m.installedRancherVersion.WithLabelValues(installedVersion).Set(1)

	elapsed := time.Since(start)
	log.Debugf("getInstalledRancherVersion metric collection took %s", elapsed)
}

func getClusterConnectedState(client rancher.Client, m rancherMetrics) {
	state, err := client.GetClusterConnectedState()
	if err != nil {
		log.Errorf("error retrieving cluster connected states: %v", err)
	}
	for key, value := range state {
		if value == true {
			m.clusterConditionConnected.WithLabelValues(key).Set(1)
			m.clusterConditionNotConnected.WithLabelValues(key).Set(0)
		} else {
			m.clusterConditionNotConnected.WithLabelValues(key).Set(1)
			m.clusterConditionConnected.WithLabelValues(key).Set(0)
		}
	}
}

func getNumberOfClusters(client rancher.Client, m rancherMetrics) {
	numberOfClusters, err := client.GetNumberOfManagedClusters()

	if err != nil {
		log.Errorf("error retrieving number of managed clusters: %v", err)
	}
	m.managedClusterCount.Set(float64(numberOfClusters))
}

func getDistributions(client rancher.Client, m rancherMetrics) {
	distributions, err := client.GetK8sDistributions()
	if err != nil {
		log.Errorf("error retrieving cluster k8s distributions: %v", err)
	}
	m.managedRKEClusterCount.Set(float64(distributions["rke"]))
	m.managedRKE2ClusterCount.Set(float64(distributions["rke2"]))
	m.managedK3sClusterCount.Set(float64(distributions["k3s"]))
	m.managedEKSClusterCount.Set(float64(distributions["eks"]))
	m.managedAKSClusterCount.Set(float64(distributions["aks"]))
	m.managedGKEClusterCount.Set(float64(distributions["gke"]))
}

func getNumberOfNodes(client rancher.Client, m rancherMetrics) {
	numberOfNodes, err := client.GetNumberOfManagedNodes()

	if err != nil {
		log.Errorf("error retrieving number of managed nodes: %v", err)
	}

	m.managedNodeCount.Set(float64(numberOfNodes))
}

func getDownstreamClusterVersions(client rancher.Client, m rancherMetrics) {
	downstreamClusterVersions, err := client.GetDownstreamClusterVersions()

	if err != nil {
		log.Errorf("error retrieving downstream k8s cluster versions: %v", err)
	}

	for _, value := range downstreamClusterVersions {

		m.downstreamClusterVersion.WithLabelValues(value.Name, value.Version).Set(1)
	}
}

func getNumberOfUsers(client rancher.Client, m rancherMetrics) {
	users, err := client.GetNumberOfUsers()
	if err != nil {
		log.Errorf("error retrieving number of users: %v", err)
	}

	m.userCount.Set(float64(users))
}

func getNumberOfTokens(client rancher.Client, m rancherMetrics) {
	tokens, err := client.GetNumberOfTokens()
	if err != nil {
		log.Errorf("error retrieving number of tokens: %v", err)
	}

	m.tokenCount.Set(float64(tokens))
}

func getNumberOfProjects(client rancher.Client, m rancherMetrics) {
	projects, err := client.GetNumberofProjects()
	if err != nil {
		log.Errorf("error retrieving number of projects: %v", err)
	}

	m.projectCount.Set(float64(projects))
}

func getProjectLabels(client rancher.Client, m rancherMetrics) {
	projectLabels, err := client.GetProjectLabels()
	if err != nil {
		log.Errorf("error retrieving project labels: %v", err)
	}

	for _, value := range projectLabels {

		m.projectLabels.WithLabelValues(value.ProjectClusterName, value.Projectid, value.ProjectDisplayName, value.LabelKey, value.LabelValue).Set(1)

	}
}

func getProjectAnnotations(client rancher.Client, m rancherMetrics) {
	projectAnnotations, err := client.GetProjectAnnotations()
	if err != nil {
		log.Errorf("error retrieving project annotations: %v", err)
	}

	for _, value := range projectAnnotations {
		m.projectAnnotations.WithLabelValues(value.ProjectClusterName, value.Projectid, value.ProjectDisplayName, value.AnnotationKey, value.AnnotationValue).Set(1)
	}
}

func getProjectResources(client rancher.Client, m rancherMetrics) {
	projectResources, err := client.GetProjectResourceQuota()
	if err != nil {
		log.Errorf("error retrieving project resources: %v", err)
	}
	for _, value := range projectResources {
		m.projectResources.WithLabelValues(value.ProjectClusterName, value.Projectid, value.ProjectDisplayName, value.ResourceKey, value.ResourceType).Set(value.ResourceValue)
	}
}

func getRancherCustomResources(client rancher.Client, m rancherMetrics) {
	resources, err := client.GetRancherCustomResourceCount()
	if err != nil {
		return
	}
	for key, value := range resources {
		m.rancherCustomResources.WithLabelValues(key).Set(float64(value))
	}
}

func getNodeInfo(client rancher.Client, m rancherMetrics) {
	nodeResources, err := client.GetManagedNodeInfo()
	if err != nil {
		return
	}

	for _, value := range nodeResources {
		m.managedNodeInfo.WithLabelValues(value.Name, value.ParentCluster, strconv.FormatBool(value.IsControlPlane),
			strconv.FormatBool(value.IsEtcd), strconv.FormatBool(value.IsWorker), value.Architecture,
			value.ContainerRuntimeVersion, value.KernelVersion, value.OS, value.OSImage).Set(1)
	}
}

func getNumberOfBackups(client rancher.Client, m rancherBackupMetrics) {
	backups, err := client.GetNumberOfBackups()
	if err != nil {
		log.Errorf("error retrieving number of backups: %v", err)
	}

	m.backupCount.Set(float64(backups))
}

func getNumberOfRestores(client rancher.Client, m rancherBackupMetrics) {
	restores, err := client.GetNumberOfRestores()
	if err != nil {
		log.Errorf("error retrieving number of restores: %v", err)
	}

	m.restoreCount.Set(float64(restores))
}

func getBackups(client rancher.Client, m rancherBackupMetrics) {
	backups, err := client.GetBackups()
	if err != nil {
		log.Errorf("error retrieving backups: %v", err)
	}

	for _, value := range backups {
		m.backup.WithLabelValues(value.Name, value.ResourceSetName, strconv.FormatInt(value.RetentionCount, 10), value.BackupType, value.Message, value.Filename, value.StorageLocation, value.NextSnapshot, value.LastSnapshot).Set(1)
	}
}

func getRestores(client rancher.Client, m rancherBackupMetrics) {
	restores, err := client.GetRestores()
	if err != nil {
		log.Errorf("error retrieving backups: %v", err)
	}

	for _, value := range restores {
		m.restore.WithLabelValues(value.Name, value.Filename, strconv.FormatBool(value.Prune), value.StorageLocation, value.Message, value.ResoreCompletionTime).Set(1)
	}
}

func getDownstreamClusterStates(client rancher.Client, m rancherMetrics) {
	rancherServerURL, err := client.GetRancherServerUrl()
	if err != nil {
		log.Errorf("error retrieving rancher server Url: %v", err)
	}

	downstreamClustersInfo, err := client.GetDownstreamClustersInfo()
	if err != nil {
		log.Errorf("error retrieving downstream k8s cluster info: %v", err)
	}

	for _, value := range downstreamClustersInfo {
		m.downstreamClusterStatus.WithLabelValues(value.Name, value.DisplayName, value.Version, rancherServerURL).Set(1)
	}
}

func getClusterConditions(client rancher.Client, m rancherMetrics) {
	clustersConditions, err := client.GetClusterConditions()
	if err != nil {
		log.Errorf("error retrieving cluster connected states: %v", err)
	}

	for key, value := range clustersConditions {
		m.clusterConditionPending.WithLabelValues(key).Set(0)
		if value["Pending"].Status == true {
			m.clusterConditionPending.WithLabelValues(key).Set(1)
		}

		m.clusterConditionWaiting.WithLabelValues(key).Set(0)
		if value["Waiting"].Status == true {
			m.clusterConditionWaiting.WithLabelValues(key).Set(1)
		}

		m.clusterConditionDiskPressure.WithLabelValues(key).Set(0)
		if value["DiskPressure"].Status == true {
			m.clusterConditionDiskPressure.WithLabelValues(key).Set(1)
		}

		m.clusterConditionMemoryPressure.WithLabelValues(key).Set(0)
		if value["MemoryPressure"].Status == true {
			m.clusterConditionMemoryPressure.WithLabelValues(key).Set(1)
		}

		m.clusterConditionReady.WithLabelValues(key, value["Ready"].Reason, value["Ready"].Message).Set(0)
		if value["Ready"].Status == true {
			m.clusterConditionReady.WithLabelValues(key, value["Ready"].Reason, value["Ready"].Message).Set(1)
		}
	}
}

// Reset GaugeVecs on each tick - facilitate state transition
func resetGaugeVecMetrics(m rancherMetrics) {
	m.installedRancherVersion.Reset()
	m.clusterConditionConnected.Reset()
	m.clusterConditionNotConnected.Reset()
	m.clusterConditionPending.Reset()
	m.clusterConditionWaiting.Reset()
	m.clusterConditionDiskPressure.Reset()
	m.clusterConditionMemoryPressure.Reset()
	m.clusterConditionReady.Reset()
	m.downstreamClusterVersion.Reset()
	m.downstreamClusterStatus.Reset()
	m.projectLabels.Reset()
	m.projectAnnotations.Reset()
	m.projectResources.Reset()
	m.managedNodeInfo.Reset()
}

func resetBackupGaugeVecMetrics(m rancherBackupMetrics) {
	m.backup.Reset()
	m.restore.Reset()
}
