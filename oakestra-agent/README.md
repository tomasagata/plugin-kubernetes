# Oakestra Agent

This ReadMe explains the functionality of the Oakestra Agent, which is responsible for establishing communication between Oakestra Root and the Kubernetes Cluster, as well as registering the cluster.


1. [Deployment](#oakestra-agent)
2. [Functionality](#Functionality)

## Deployment
For the deployment of the component, a Kubernetes deployment simply needs to be initiated. However, it is important that the necessary environment variables are set. 


| Variable Name                   | Default Values      | Description                                                     |
|---------------------------------|---------------------|-----------------------------------------------------------------|
| ROOT_SYSTEM_MANAGER_IP          | 192.168.123.225     | IP Oakestra Root                                                |
| ROOT_SYSTEM_MANAGER_PORT        | 10000               | Port Oakestra Root                                              |
| ROOT_SERVICE_MANAGER_PORT       | 10099               | Port Oakestra Network Root                                      |
| ROOT_GRPC_PORT                  | 50052               | Port GRPC Root                                                  |
| CLUSTER_NAME                    | *Needs to be set*  | Name of Cluster                                                 |
| CLUSTER_LOCATION                | *Needs to be set*  | Location of Cluster                                             |
| MY_PORT                         | 10100               | Local port which starts server                                  |
| NODE_PORT                       | 30000               | Exposed public port to Root, needs to be in range 30000-32767   |
| CLUSTER_SERVICE_MANAGER_PORT    | 30330               | NodePort for Cluster Service Manager                            |
| CLUSTER_SERVICE_MANAGER_IP      | localhost           | Node IP for Cluster Service Manager                             |



## Functionality

The Oakestra Kubernetes Agent is a Kubernetes extension implemented in Go that mirrors the endpoints of an Oakestra Cluster Orchestrator. This extension allows seamless communication between Oakestra Root and Kubernetes clusters, enabling interaction with Kubernetes resources through familiar Oakestra endpoints.

### Client Registration 

The directory [clusterRegistration](./agent/clusterRegistration/) implements a client in Go that registers itself with Oakestra Root using gRPC. The client communicates with the Oakestra Root server to handle the initialization and finalization processes.

Update the addr variable with the correct address of the Oakestra Root server.
The client establishes a gRPC connection to the Oakestra Root server and invokes the following methods:
HandleInitGreeting: Initiates the registration process by sending a greeting message.
HandleInitFinal: Completes the registration process by providing additional information, such as manager and network component ports, cluster name, cluster location, and custom cluster information.



### Hardware Aggregation

Thes directory [aggregation](./agent/aggregation) provides an integration between a Kubernetes cluster and Oakestra Root. It collects information about the cluster's resource usage and sends it to an Oakestra Root server.

The program aggregates information about the Kubernetes cluster, including CPU and memory usage, GPU information (if available), and general node details. This information is then formatted into JSON and sent to the specified Oakestra Root server.


### Kubernetes Client

The [kubenetesClient](./agent/kubernetesClient) package provides functionalities for interacting with the Oakestra Kubernetes Controller.

