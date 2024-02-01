package rancher

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	settingGVR = schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "settings"}
	clusterGVR = schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "clusters"}
	nodeGVR    = schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "nodes"}
	tokenGVR   = schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "tokens"}
	usersGVR   = schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "users"}
)

type Client struct {
	Client dynamic.Interface
	Config *rest.Config
}

type clusterVersion struct {
	Name    string
	Version string
}

type clusterInfo struct {
	Name        string
	DisplayName string
	Version     string
}

type projectLabel struct {
	Projectid          string
	ProjectDisplayName string
	ProjectClusterName string
	LabelKey           string
	LabelValue         string
}

type projectAnnotation struct {
	Projectid          string
	ProjectDisplayName string
	ProjectClusterName string
	AnnotationKey      string
	AnnotationValue    string
}

type projectResource struct {
	Projectid          string
	ProjectDisplayName string
	ProjectClusterName string
	ResourceKey        string
	ResourceValue      float64
	ResourceType       string
}

type Release struct {
	TagName string `json:"tag_name"`
}

type nodeInfo struct {
	Name                    string
	ParentCluster           string
	IsControlPlane          bool
	IsEtcd                  bool
	IsWorker                bool
	Architecture            string
	ContainerRuntimeVersion string
	KernelVersion           string
	OS                      string
	OSImage                 string
}

func (r Client) GetInstalledRancherVersion() (string, error) {

	res, err := r.Client.Resource(settingGVR).Get(context.Background(), "server-version", v1.GetOptions{})
	if err != nil {
		return "", err
	}

	version, _, err := unstructured.NestedString(res.UnstructuredContent(), "value")
	if err != nil {
		return "", err
	}
	return version, nil
}

func (r Client) GetNumberOfManagedClusters() (int, error) {

	res, err := r.Client.Resource(clusterGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return 0, err
	}

	return len(res.Items), nil
}

func (r Client) GetK8sDistributions() (map[string]int, error) {

	distributions := make(map[string]int)

	res, err := r.Client.Resource(clusterGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, v := range res.Items {
		labels := v.GetLabels()
		distribution := labels["provider.cattle.io"]
		distributions[distribution] += 1
	}

	return distributions, nil
}

func (r Client) GetLatestRancherVersion() (string, error) {

	// check dns resultion first, if this fails http.Get segfaults
	_, err := net.LookupHost("api.github.com")
	if err != nil {
		return "", err
	}

	resp, err := http.Get("https://api.github.com/repos/rancher/rancher/releases")
	if err != nil {
		fmt.Printf("failed to retrieve releases: %v\n", err)
		return "", err
	}

	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	// Parse JSON response
	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		fmt.Printf("failed to parse JSON: %v\n", err)
		return "", err
	}

	// Find latest release version that is not a pre-release/release candidate/patch version
	var latestVersion string
	re := regexp.MustCompile(`^v\d+\.\d+\.\d+$`)
	for _, release := range releases {
		if re.MatchString(release.TagName) && release.TagName > latestVersion {
			latestVersion = release.TagName
		}
	}

	return latestVersion, nil
}

func (r Client) GetNumberOfManagedNodes() (int, error) {

	res, err := r.Client.Resource(nodeGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return 0, err
	}

	return len(res.Items), nil
}

func (r Client) GetManagedNodeInfo() ([]nodeInfo, error) {
	var nodes []nodeInfo

	res, err := r.Client.Resource(nodeGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, node := range res.Items {

		nodeValue := nodeInfo{}

		parentClusterID, _, err := unstructured.NestedString(node.Object, "metadata", "namespace")
		if err != nil {
			return nil, err
		}

		clusterName, err := r.clusterIdToName(parentClusterID)
		if err != nil {
			return nil, err
		}

		if clusterName != "local" {

			nodeValue.ParentCluster = clusterName

			name, _, err := unstructured.NestedString(node.Object, "spec", "requestedHostname")
			if err != nil {
				return nil, err
			}

			nodeValue.Name = name

			cpl, _, err := unstructured.NestedBool(node.Object, "spec", "controlPlane")
			if err != nil {
				return nil, err
			}

			nodeValue.IsControlPlane = cpl

			etcd, _, err := unstructured.NestedBool(node.Object, "spec", "etcd")
			if err != nil {
				return nil, err
			}

			nodeValue.IsEtcd = etcd

			worker, _, err := unstructured.NestedBool(node.Object, "spec", "worker")
			if err != nil {
				return nil, err
			}

			nodeValue.IsWorker = worker

			architecture, _, err := unstructured.NestedString(node.Object, "status", "internalNodeStatus", "nodeInfo", "architecture")
			if err != nil {
				return nil, err
			}

			nodeValue.Architecture = architecture

			containerRunTime, _, err := unstructured.NestedString(node.Object, "status", "internalNodeStatus", "nodeInfo", "containerRuntimeVersion")
			if err != nil {
				return nil, err
			}

			nodeValue.ContainerRuntimeVersion = containerRunTime

			kernelVersion, _, err := unstructured.NestedString(node.Object, "status", "internalNodeStatus", "nodeInfo", "kernelVersion")
			if err != nil {
				return nil, err
			}

			nodeValue.KernelVersion = kernelVersion

			OS, _, err := unstructured.NestedString(node.Object, "status", "internalNodeStatus", "nodeInfo", "operatingSystem")
			if err != nil {
				return nil, err
			}

			nodeValue.OS = OS

			osImage, _, err := unstructured.NestedString(node.Object, "status", "internalNodeStatus", "nodeInfo", "osImage")
			if err != nil {
				return nil, err
			}

			nodeValue.OSImage = osImage

			nodes = append(nodes, nodeValue)
		}
	}

	return nodes, nil
}

func (r Client) GetClusterConnectedState() (map[string]bool, error) {
	clusterStatus := make(map[string]bool)

	res, err := r.Client.Resource(clusterGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Iterate through each cluster management object
	for _, cluster := range res.Items {

		// Grab Cluster Name
		clusterName, _, err := unstructured.NestedString(cluster.Object, "spec", "displayName")

		if err != nil {
			return nil, err
		}

		// Ignore if it's the "Local" cluster
		if clusterName != "local" {

			clusterStatus[clusterName] = false

			// Grab status.conditions slice from object
			statusSlice, _, _ := unstructured.NestedSlice(cluster.Object, "status", "conditions")

			// Iterate through each status slice to determine if cluster is connected
			for _, value := range statusSlice {

				// Determine whether we find both conditions in each status Message
				// We're looking for both when type == connected and status == true
				// to identify if a cluster is connected to this Rancher instance

				// Reset values when iterating through
				foundStatus := false
				foundType := false

				for k, v := range value.(map[string]interface{}) {

					if k == "type" && v.(string) == "Connected" {
						foundType = true
					}

					if k == "status" && v.(string) == "True" {
						foundStatus = true
					}

					if foundStatus == true && foundType == true {
						clusterStatus[clusterName] = true
					}

				}
			}

		}
	}

	return clusterStatus, nil
}

func (r Client) GetDownstreamClusterVersions() ([]clusterVersion, error) {

	var clusters []clusterVersion

	res, err := r.Client.Resource(clusterGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Iterate through each cluster management object
	for _, cluster := range res.Items {

		// Grab Cluster Name
		clusterName, _, err := unstructured.NestedString(cluster.Object, "spec", "displayName")

		if err != nil {
			return nil, err
		}

		clusterK8sVersion, _, err := unstructured.NestedString(cluster.Object, "status", "version", "gitVersion")

		if err != nil {
			return nil, err
		}

		if clusterK8sVersion != "" {

			clusterInfo := clusterVersion{
				Name:    clusterName,
				Version: clusterK8sVersion,
			}

			clusters = append(clusters, clusterInfo)
		}
	}

	return clusters, nil

}

func (r Client) GetNumberOfTokens() (int, error) {

	res, err := r.Client.Resource(tokenGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return 0, err
	}

	return len(res.Items), nil
}

func (r Client) GetNumberOfUsers() (int, error) {

	res, err := r.Client.Resource(usersGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return 0, err
	}

	var users int

	for _, user := range res.Items {
		userName, _, err := unstructured.NestedString(user.Object, "username")
		if err != nil {
			return 0, err
		}
		if userName != "" {
			users++
		}
	}

	return users, nil
}

func (r Client) GetRancherServerUrl() (string, error) {

	res, err := r.Client.Resource(settingGVR).Get(context.Background(), "server-url", v1.GetOptions{})
	if err != nil {
		return "", err
	}

	url, _, err := unstructured.NestedString(res.UnstructuredContent(), "value")
	if err != nil {
		return "", err
	}
	return url, nil
}

func (r Client) GetDownstreamClustersInfo() ([]clusterInfo, error) {

	var clusterInfos []clusterInfo

	res, err := r.Client.Resource(clusterGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Iterate through each cluster management object
	for _, cluster := range res.Items {

		// Grab Cluster Name
		displayName, _, err := unstructured.NestedString(cluster.Object, "spec", "displayName")
		name, _, err := unstructured.NestedString(cluster.Object, "metadata", "name")

		if err != nil {
			return nil, err
		}

		clusterK8sVersion, _, err := unstructured.NestedString(cluster.Object, "status", "version", "gitVersion")

		if err != nil {
			return nil, err
		}

		clusterInfo := clusterInfo{
			Name:        name,
			DisplayName: displayName,
			Version:     clusterK8sVersion,
		}

		clusterInfos = append(clusterInfos, clusterInfo)
	}

	return clusterInfos, nil

}

func (r Client) GetClusterConditions() (map[string]map[string]bool, error) {
	clusterConditions := make(map[string]map[string]bool)

	res, err := r.Client.Resource(clusterGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Iterate through each cluster management object
	for _, cluster := range res.Items {

		// Grab Cluster Name
		clusterName, _, err := unstructured.NestedString(cluster.Object, "metadata", "name")

		if err != nil {
			return nil, err
		}

		// Ignore if it's the "Local" cluster
		if clusterName != "local" {

			clusterConditions[clusterName] = make(map[string]bool)

			// Grab status.conditions slice from object
			conditions, _, _ := unstructured.NestedSlice(cluster.Object, "status", "conditions")

			// Iterate through each status slice to determine if cluster is connected
			for _, value := range conditions {

				// Determine whether we find both conditions in each status Message
				// We're looking for both when type == connected and status == true
				// to identify if a cluster is connected to this Rancher instance

				// Reset values when iterating through
				foundStatus := false
				foundType := ""

				for k, v := range value.(map[string]interface{}) {

					if k == "type" && v.(string) == "Pending" {
						foundType = "Pending"
					}

					if k == "status" {
						foundStatus = v.(string) == "True"
					}

					if k == "type" && v.(string) == "Waiting" {
						foundType = "Waiting"
					}

					if k == "status" {
						foundStatus = v.(string) == "True"
					}

					if k == "type" && v.(string) == "DiskPressure" {
						foundType = "DiskPressure"
					}

					if k == "status" {
						foundStatus = v.(string) == "True"
					}

					if k == "type" && v.(string) == "MemoryPressure" {
						foundType = "MemoryPressure"
					}

					if k == "status" {
						foundStatus = v.(string) == "True"
					}

					if k == "type" && v.(string) == "Ready" {
						foundType = "Ready"
					}

					if k == "status" {
						foundStatus = v.(string) == "True"
					}

					if foundType != "" {
						clusterConditions[clusterName][foundType] = foundStatus
					}

				}
			}

		}
	}

	return clusterConditions, nil
}
