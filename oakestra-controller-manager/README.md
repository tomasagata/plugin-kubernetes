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

```sh
kubectl apply -f dist/install
```

### To remove from the cluster

```sh
kubectl delete -f dist/install
```

### To develop new features or changes

Build the installer with
```sh
# Replace the image with the image name of your choice.
# If you don't want to keep adding the IMG environment variable, you can change it
# from the Makefile
make build-installer IMG=<Your image name of choice>
```

Then build the image using
```sh
make docker-buildx IMG=<Your image name of choice>
```


> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin 
privileges or be logged in as admin.



