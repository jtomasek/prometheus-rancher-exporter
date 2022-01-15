# Prometheus Exporter for Rancher

# Current Metrics Scraped:

* Rancher Major/Minor/Patch versions
* Total number of Managed clusters
* Number of RKE/RKE2/K3s/EKS/AKS/GKE clusters

# Prereqs

* Decide how to auth with the `local` cluster:

## In-cluster config

Uncomment the following from `main.go`

```go
// Use this for in-cluster config 
//config, err := rest.InClusterConfig()
```

## External cluster config

(default) and is handled by the following code in `main.go`

```go
// Use this for out of cluster config
currentUser, err := user.Current()
if err != nil {
	log.Fatal(err.Error())
}

kubeconfig := flag.String("kubeconfig", fmt.Sprintf("/home/%s/.kube/config", currentUser.Username), "absolute path to the kubeconfig file")
flag.Parse()
config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
```

# Testing

`go run main.go` and access `http://localhost:8080`