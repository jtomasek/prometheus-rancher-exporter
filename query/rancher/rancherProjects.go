package rancher

/*
Number of Projects
Projects to Namespace
Project Resources

Examples:

rancher_project_labels{project_id="p-6rv9k",project_display_name="xy",label_appteam="xy-team"} 1
rancher_project_annotations{project_id="p-6rv9k",project_display_name="xy",annotation_custom="bla"} 1

rancher_namespace_info{project_id="p-6rv9k",project_display_name="xy",namespace_display_name="my_ns1", namespace_id="..."} 1

kube_resourcequota{namespace="my_ns1", resource="configmaps",resourcequota="default-8kt7x",type="hard"}	 30
kube_resourcequota{namespace="my_ns1", resource="limits.cpu",resourcequota="default-8kt7x",type="hard"}	 26
kube_resourcequota{namespace="my_ns1", resource="limits.memory",resourcequota="default-8kt7x",type="hard"}	 137438953472
kube_resourcequota{namespace="my_ns1", resource="persistentvolumeclaims",resourcequota="default-8kt7x",type="hard"}	 150
kube_resourcequota{namespace="my_ns1", resource="pods",resourcequota="default-8kt7x",type="hard"}	 100
kube_resourcequota{namespace="my_ns1", resource="replicationcontrollers",resourcequota="default-8kt7x",type="hard"}	 100
kube_resourcequota{namespace="my_ns1", resource="requests.cpu",resourcequota="default-8kt7x",type="hard"}	 26
kube_resourcequota{namespace="my_ns1", resource="requests.memory",resourcequota="default-8kt7x",type="hard"}	 137438953472
kube_resourcequota{namespace="my_ns1", resource="requests.storage",resourcequota="default-8kt7x",type="hard"}	 549755813888
kube_resourcequota{namespace="my_ns1", resource="secrets",resourcequota="default-8kt7x",type="hard"}	 130
kube_resourcequota{namespace="my_ns1", resource="services",resourcequota="default-8kt7x",type="hard"}	 30
kube_resourcequota{namespace="my_ns1", resource="services.loadbalancers",resourcequota="default-8kt7x",type="hard"}	 1
kube_resourcequota{namespace="my_ns1", resource="services.nodeports",resourcequota="default-8kt7x",",type="hard"}	 30

*/

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	projectsGVR = schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "projects"}
)

func (r Client) GetNumberofProjects() (int, error) {
	res, err := r.Client.Resource(projectsGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return 0, err
	}

	return len(res.Items), nil
}

func (r Client) GetProjectLabels() ([]projectLabel, error) {
	res, err := r.Client.Resource(projectsGVR).List(context.Background(), v1.ListOptions{})

	if err != nil {
		return nil, err
	}

	var projectLabelsArray []projectLabel

	// Loop through array of Projects
	for _, projectValue := range res.Items {

		projectLabels := projectValue.GetLabels()

		projectDisplayName, _, err := unstructured.NestedString(projectValue.Object, "spec", "displayName")
		if err != nil {
			return nil, err
		}

		fmt.Println("test")

		// Loop through array of Labels within each Project
		for labelKey, labelValue := range projectLabels {

			fmt.Printf("Project Name: %s, Label Key: %s, Label Value: %s", projectDisplayName, labelKey, labelValue)
			fmt.Println("Test2")
		}
	}
	return projectLabelsArray, nil
}

// Projects return the cluster ID (ie c-m-xwf4csvg). Helper function used to lookup the display name
func (r Client) clusterIdToName(id string) (string, error) {

	res, err := r.Client.Resource(clusterGVR).Get(context.Background(), id, v1.GetOptions{})
	if err != nil {
		return "", err
	}

	clusterName, _, err := unstructured.NestedString(res.Object, "spec", "displayName")
	if err != nil {
		return "", err
	}

	return clusterName, nil

}
