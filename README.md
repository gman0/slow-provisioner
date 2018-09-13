# slow-provisioner

slow-provisioner is a Kubernetes external-storage plugin. It fulfills persistent volume claims by provisioning *fake* volumes (not represented by any actual storage) which are then meant to be consumed by the [slow-csi](https://github.com/gman0/slow-csi) CSI plugin. It simulates pending operations with configurable delays, useful for debugging, triggering timeouts, etc.

Delayable operations:
* `Provision`
* `Delete`

## Building

slow-provisioner can be compiled in a form of a binary file or in a form of a Docker image. When compiled as a binary file, the result is stored in `_output/` directory with the name `slow-provisioner`. When compiled as an image, it's stored in the local Docker image store.

Building binary:
```bash
$ make provisioner
```

Building Docker image:
```bash
$ make image
```

## Configuration

**Available command line arguments:**

Option | Default value | Description
------ | ------------- | -----------
`--kubeconfig` | _none_ | Path to a kube config. Only required if out-of-cluster
`--provisioner` | `slow-provisioner` | Name of the provisioner. The provisioner will only provision volumes for claims that request a StorageClass with a provisioner field set equal to this name
`--nodeplugin` | `csi-slowplugin` | Name of the CSI node plugin
`--defaultdelay` | _none_ | Default delay applied to all provisioner operations
`--delay` | _none_ | Delay settings for individual provisioner operations (e.g. `Provision=20..50,Delete=inf`)

**Available StorageClass parameters:**  
None.

### Delay format

All delays are in seconds and may be in one of three formats:
* `n`: simple number
* `a..b`: random delay in an interval _[`a`, `b`)_, `b` must be greater than `a`
* `inf`: infinite delay, such RPC will never finish

Example: `--delay=CreateVolume=10..20,DeleteVolume=5,NodeUnstageVolume=inf`

## Deployment

Requires Kubernetes 1.11+

YAML manifests are located in `deploy/kubernetes/`.

```bash
$ kubectl create -f rbac.yaml
$ kubectl create -f deployment.yaml
```

Deploys RBACs required for the provisioner to run, and the slow-provisioner deployment itself where you can configure the delays.

Manifests with example StorageClass, PVC and a deployment are located in `exapmle/`:

```bash
$ kubectl create -f storageclass.yaml
$ kubectl create -f pvc.yaml
$ kubectl create -f deployment.yaml
```
