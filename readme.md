![workflow status](https://github.com/David-VTUK/prometheus-rancher-exporter/actions/workflows/test-build-publish.yml/badge.svg)

# Unofficial Prometheus Exporter for Rancher

**Note** : This project is not officially supported by Rancher/Suse.

# Quickstart

1. Enable monitoring in the Rancher Management Cluster, aka `local` cluster
2. Apply the manifest from this repo : `kubectl apply -f ./manifests/exporter.yaml`

# Grafana Dashboard

`./manifests/grafana-dashboard` includes a basic dashboard in JSON format that can be imported into Grafana.  
`./manifests/grafana-dashboard-projects.json` includes a Rancher project-focused dashboard in JSON format that can be imported into Grafana.  
`./manifests/grafana-dashboard-all-cr.json` includes a dynamic dashboard showing counts for all Rancher custom resources (*.cattle.io).  


![img.png](img/overview-dashboard.png)
![img.png](img/proj-dashboard.png)
![img.png](img/cr-dashboard.png)

# Developing

By default, the exporter will use in-cluster authentication via a associated service account

## External cluster config

To test using external authentication via the local `kubeconfig`, set the following environment variable:

```bash
export RANCHER_EXPORTER_EXTERNAL_AUTH=TRUE
```

* `go run main.go`
* Navigate to http://localhost:8080/metrics
