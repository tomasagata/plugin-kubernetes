# plugin-kubernetes

This is the repository for integrating Kubernetes with Oakestra. A Kubernetes cluster can be added to Oakestra as an additional Oakestra cluster. The following setup must be pursued for this integration.

All these components must be installed on the Kubernetes cluster.

   

- [plugin-kubernetes](#plugin-kubernetes)
    - [0. Prerequisites](#0-prerequisites)
      - [1.1 Create Namespaces in Kubernetes](#11-create-namespaces-in-kubernetes)
    - [2. Oakestra Network](#2-oakestra-network)
      - [1.1 Oakestra CNI](#11-oakestra-cni)
      - [1.2 Multus CNI](#12-multus-cni)
      - [1.3 Oakestra Cluster Service Manager](#13-oakestra-cluster-service-manager)
      - [1.4 Oakestra Node NetManager](#14-oakestra-node-netmanager)
    - [2. Oakestra Kubernetes Agent](#2-oakestra-kubernetes-agent)
    - [3. Oakestra Kubernetes Controller](#3-oakestra-kubernetes-controller)



*Important Note:*\
**Many of the Kubernetes resources require environment variables to be set. For more details, refer to the respective READMEs of the components. The components need to be initiated in the sequence they are listed here.**

### 0. Prerequisites
For the integration of Kubernetes with Oakestra, a few prerequisites must be met beforehand. Firstly, there must be an existing Kubernetes cluster with kubectl access. Secondly, all nodes of the Kubernetes cluster must be able to communicate with all nodes of [Oakestra](https://github.com/oakestra). Moreover, a default CNI (e.g., Calico) is required. Lastly, an Oakestra Root server must be operational to facilitate this integration.

#### 1.1 Create Namespaces in Kubernetes

```bash
kubectl create namespace oakestra-system

kubectl create namespace oakestra

kubectl create namespace oakestra-controller-manager
```


### 2. Oakestra Network

#### 1.1 Oakestra CNI
To communicate with additional Oakestra resources, a separate Container Network Interface (CNI) is required, which must be installed on all nodes. For this purpose, a DaemonSet is used to ensure that all necessary installations are automatically performed by the NetManager on each node in the cluster. This DeamonSet is deployed in Chapter [3.4 Oakestra Node NetManager](#34-oakestra-node-netmanager).

#### 1.2 Multus CNI
To use two CNIs per container, Multus is required. The following commands must be executed:


```bash
kubectl apply -f https://raw.githubusercontent.com/k8snetworkplumbingwg/multus-cni/master/deployments/multus-daemonset.yml

kubectl apply -f oakestra-network/Deployment/oakestra-cni/oakesta-cni.yaml -n oakestra
```


#### 1.3 Oakestra Cluster Service Manager
This component needs to run once per cluster and requires a MongoDB and an MQTT server.


```bash
kubectl apply -f oakestra-network/Deployment/oakestra-cluster-service-manager/mosquitto/ -n oakestra-system

kubectl apply -f oakestra-network/Deployment/oakestra-cluster-service-manager/mongodb/ -n oakestra-system

```

**Check ReadMe for Oakestra Cluster Service Manager**
```bash
kubectl apply -f oakestra-network/Deployment/oakestra-cluster-service-manager/oakestra-cluster-service-manager.yaml -n oakestra-system
```

#### 1.4 Oakestra Node NetManager
This component must also run on all nodes; it is responsible for ensuring that the containers in Kubernetes find the correct routing.

```bash

kubectl set env daemonset/calico-node -n kube-system IP_AUTODETECTION_METHOD="skip-interface=goProxy.*"

kubectl apply -f oakestra-network/Deployment/oakestra-nodenetmanager/node-netmanager.yaml -n oakestra-system
```


### 2. Oakestra Kubernetes Agent
The agent is needed to register with the cluster and establish communication with Root.

The following Kubernetes deployment must be initiated:


```bash
kubectl apply -f oakestra-agent/Deployment/oakestra-agent.yaml
```


### 3. Oakestra Kubernetes Controller
This controller-manager is responsible for deploying the appropriate resources in Kubernetes for Oakestra.

The following commands must be executed:


**Certmanager must be installed for this.** 
The current documentation can be found [here](https://cert-manager.io/docs/installation/). 


```bash
cd oakestra-controller-manager

make install

//TODO  Not yet in oakestra cr.
make deploy IMG=ghcr.io/oakestra/oakestra-controller:1.0 
```