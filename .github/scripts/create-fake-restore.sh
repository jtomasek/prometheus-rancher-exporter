#!/usr/bin/bash

kubectl config use-context k3d-upstream
kubectl apply -f - <<EOF
apiVersion: resources.cattle.io/v1
kind: Restore
metadata:
  creationTimestamp: '2023-07-19T11:20:25Z'
  generateName: restore-
  generation: 1
  managedFields:
    - apiVersion: resources.cattle.io/v1
      fieldsType: FieldsV1
      fieldsV1:
        f:metadata:
          f:generateName: {}
        f:spec:
          .: {}
          f:backupFilename: {}
          f:deleteTimeoutSeconds: {}
          f:prune: {}
      manager: rancher
      operation: Update
      time: '2023-07-19T11:20:25Z'
    - apiVersion: resources.cattle.io/v1
      fieldsType: FieldsV1
      fieldsV1:
        f:status:
          .: {}
          f:backupSource: {}
          f:conditions: {}
          f:observedGeneration: {}
          f:restoreCompletionTs: {}
          f:summary: {}
      manager: backup-restore-operator
      operation: Update
      subresource: status
      time: '2023-07-19T11:22:07Z'
  name: restore-jq9bs
  resourceVersion: '28180618'
  uid: 87e3dbc7-9e7d-4149-b2d0-de99e143610f
spec:
  backupFilename: >-
    one-time-test-2-e3acb0dc-c4f1-4482-83db-66f0141722de-2023-07-19T11-16-41Z.tar.gz
  deleteTimeoutSeconds: 10
  prune: true
status:
  backupSource: PV
  conditions:
    - lastUpdateTime: '2023-07-19T11:22:07Z'
      message: Completed
      status: 'True'
      type: Ready
  observedGeneration: 1
  restoreCompletionTs: '2023-07-19T11:22:07Z'
  summary: ''
EOF