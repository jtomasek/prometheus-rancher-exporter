package rancher

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"regexp"
	"strconv"
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

		projectClusterID, _, err := unstructured.NestedString(projectValue.Object, "spec", "clusterName")
		if err != nil {
			return nil, err
		}

		projectClusterName, _ := r.clusterIdToName(projectClusterID)

		if projectClusterName != "" {

			for labelKey, labelValue := range projectLabels {

				project := projectLabel{
					Projectid:          projectValue.GetName(),
					ProjectDisplayName: projectDisplayName,
					ProjectClusterName: projectClusterName,
					LabelKey:           labelKey,
					LabelValue:         labelValue,
				}

				projectLabelsArray = append(projectLabelsArray, project)

			}
		}
	}
	return projectLabelsArray, nil
}

func (r Client) GetProjectAnnotations() ([]projectAnnotation, error) {
	res, err := r.Client.Resource(projectsGVR).List(context.Background(), v1.ListOptions{})

	if err != nil {
		return nil, err
	}

	var projectAnnotationsArray []projectAnnotation

	// Loop through array of Projects
	for _, projectValue := range res.Items {

		projectDisplayName, _, err := unstructured.NestedString(projectValue.Object, "spec", "displayName")
		if err != nil {
			return nil, err
		}

		projectClusterID, _, err := unstructured.NestedString(projectValue.Object, "spec", "clusterName")
		if err != nil {
			return nil, err
		}

		projectClusterName, _ := r.clusterIdToName(projectClusterID)

		projectAnnotations := projectValue.GetAnnotations()

		if projectClusterName != "" {

			for annotationKey, annotationValue := range projectAnnotations {

				annotation := projectAnnotation{
					Projectid:          projectValue.GetName(),
					ProjectDisplayName: projectDisplayName,
					ProjectClusterName: projectClusterName,
					AnnotationKey:      annotationKey,
					AnnotationValue:    annotationValue,
				}

				projectAnnotationsArray = append(projectAnnotationsArray, annotation)

			}
		}
	}
	return projectAnnotationsArray, nil
}

func (r Client) GetProjectResourceQuota() ([]projectResource, error) {

	res, err := r.Client.Resource(projectsGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var projectResourceArray []projectResource

	// Loop through array of Projects
	for _, projectValue := range res.Items {

		projectDisplayName, _, err := unstructured.NestedString(projectValue.Object, "spec", "displayName")
		if err != nil {
			return nil, err
		}

		projectClusterID, _, err := unstructured.NestedString(projectValue.Object, "spec", "clusterName")
		if err != nil {
			return nil, err
		}

		projectClusterName, _ := r.clusterIdToName(projectClusterID)

		projectResourceQuotas, _, err := unstructured.NestedMap(projectValue.Object, "spec", "resourceQuota", "limit")

		if err != nil {
			return nil, err
		}

		if projectClusterName != "" {

			for key, value := range projectResourceQuotas {

				//Strip any non-numeric values from string
				re := regexp.MustCompile("[0-9]+")
				strippedString := re.FindAllString(value.(string), -1)

				convertedValue, err := strconv.ParseFloat(strippedString[0], 64)
				if err != nil {
					return nil, err
				}

				resource := projectResource{
					Projectid:          projectValue.GetName(),
					ProjectDisplayName: projectDisplayName,
					ProjectClusterName: projectClusterName,
					ResourceKey:        key,
					ResourceValue:      convertedValue,
					ResourceType:       "hard",
				}

				projectResourceArray = append(projectResourceArray, resource)

			}

			projectResourceQuotas, _, err := unstructured.NestedMap(projectValue.Object, "spec", "resourceQuota", "usedLimit")

			if err != nil {
				return nil, err
			}

			for key, value := range projectResourceQuotas {

				//Strip any non-numeric values from string
				re := regexp.MustCompile("[0-9]+")
				strippedString := re.FindAllString(value.(string), -1)

				convertedValue, err := strconv.ParseFloat(strippedString[0], 64)
				if err != nil {
					return nil, err
				}
				resource := projectResource{
					Projectid:          projectValue.GetName(),
					ProjectDisplayName: projectDisplayName,
					ProjectClusterName: projectClusterName,
					ResourceKey:        key,
					ResourceValue:      convertedValue,
					ResourceType:       "used",
				}

				projectResourceArray = append(projectResourceArray, resource)

			}
		}
	}

	return projectResourceArray, err
}

// Projects return the cluster ID (ie c-m-xwf4csvg). Helper function used to lookup the display name
// as well as it's referencing a valid, existing cluster
func (r Client) clusterIdToName(id string) (string, error) {

	var clusterName string

	res, err := r.Client.Resource(clusterGVR).Get(context.Background(), id, v1.GetOptions{})

	if err != nil {
		return "", err
	}

	// Ensure Project is referencing a valid, existing cluster
	if res != nil {

		clusterName, _, err = unstructured.NestedString(res.Object, "spec", "displayName")
		if err != nil {
			return "", err
		}

	}
	return clusterName, nil

}
