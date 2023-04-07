#!/usr/bin/bash

kubectl config use-context k3d-upstream
kubectl apply -f - <<EOF
apiVersion: management.cattle.io/v3
kind: Project
metadata:
  labels:
    cattle.io/creator: norman
  namespace: local
  name: ci-project
spec:
  clusterName: local
  displayName: ci-project
  namespaceDefaultResourceQuota:
    limit:
      configMaps: '10'
  resourceQuota:
    limit:
      configMaps: '1000'
    usedLimit: {}
EOF