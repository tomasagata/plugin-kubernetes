# Oakestra Controller Manager

Controller for All Oakestra Resources in Kubernetes

## Description
All Oakestra resources, such as Oakestra Jobs, are represented in Kubernetes using Kubernetes resources. To ensure the correct Kubernetes resources are selected, each Oakestra resource incorporates an Operator pattern, which is housed in this repository. This includes a CRD (Custom Resource Definition) and a controller.
The Controller Manager simply needs to be initiated. The resources are then utilized by the Oakestra Agent.


## Deployment

### Prerequisites
- go version v1.20.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=ghcr.io/jakobke/oakestra-controller:1.0
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin 
privileges or be logged in as admin.



