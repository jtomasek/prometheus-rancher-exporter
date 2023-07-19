#!/usr/bin/bash

kubectl config use-context k3d-upstream
kubectl apply -f - <<EOF
apiVersion: resources.cattle.io/v1
kind: Backup
metadata:
  creationTimestamp: '2023-07-15T06:52:44Z'
  generation: 1
  managedFields:
    - apiVersion: resources.cattle.io/v1
      fieldsType: FieldsV1
      fieldsV1:
        f:spec:
          .: {}
          f:resourceSetName: {}
          f:retentionCount: {}
          f:schedule: {}
      manager: rancher
      operation: Update
      time: '2023-07-15T06:52:44Z'
    - apiVersion: resources.cattle.io/v1
      fieldsType: FieldsV1
      fieldsV1:
        f:status:
          .: {}
          f:backupType: {}
          f:conditions: {}
          f:filename: {}
          f:lastSnapshotTs: {}
          f:nextSnapshotAt: {}
          f:observedGeneration: {}
          f:storageLocation: {}
          f:summary: {}
      manager: backup-restore-operator
      operation: Update
      subresource: status
      time: '2023-07-19T00:00:09Z'
  name: test-recurring
  resourceVersion: '27796019'
  uid: 990346a9-b9bc-4893-aceb-fdcab5872c65
spec:
  resourceSetName: rancher-resource-set
  retentionCount: 10
  schedule: '@midnight'
status:
  backupType: Recurring
  conditions:
    - lastUpdateTime: '2023-07-19T00:00:09Z'
      message: Completed
      status: 'True'
      type: Ready
    - lastUpdateTime: '2023-07-19T00:00:09Z'
      status: 'True'
      type: Uploaded
  filename: >-
    test-recurring-e3acb0dc-c4f1-4482-83db-66f0141722de-2023-07-19T00-00-00Z.tar.gz
  lastSnapshotTs: '2023-07-19T00:00:09Z'
  nextSnapshotAt: '2023-07-20T00:00:00Z'
  observedGeneration: 1
  storageLocation: PV
  summary: ''
EOF

