#!/usr/bin/bash

kubectl config use-context k3d-upstream
kubectl apply -f - <<EOF
apiVersion: management.cattle.io/v3
kind: Cluster
metadata:
  annotations:
  labels:
    provider.cattle.io: k3s
  name: c-m-pskdut5m
spec:
  displayName: fake-cluster
  localClusterAuthEndpoint:
    enabled: false
  windowsPreferedCluster: false
  driver: k3s
  nodeCount: 1
  provider: k3s
EOF

