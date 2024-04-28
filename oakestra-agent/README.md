# Oakestra-Kubernetes Integration

Dieses Repo beinhaltet alle wichtigen Komponenten, die für die Integration von Kubernetes in ein Oakestra Cluster notwendig sind. 

Jede einzelne Komponente ist weiter unten ausführlich erklärt, sowie das Deployment. 

Mehr Information zu Oakestra und der Infrastuktur finden sie hier:


1. [Oakestra Kubernetes Agent](#oakestra-agent)
2. [Oakestra Job Controller](#oakestra-job-controller)
3. [Oakestra CNI](#oakestra-job-controller)
4. [Oakestra Node Netmanager](#oakestra-job-controller)

## Oakestra Kubernetes Agent

The Oakestra Kubernetes Agent is a Kubernetes extension implemented in Go that mirrors the endpoints of an Oakestra Cluster Orchestrator. This extension allows seamless communication between Oakestra Root and Kubernetes clusters, enabling interaction with Kubernetes resources through familiar Oakestra endpoints.

- Diese und noch mehrere
- 


## Client Registration 

The directory [oakestra-agent](./oakestra-agent/) implements a client in Go that registers itself with Oakestra Root using gRPC. The client communicates with the Oakestra Root server to handle the initialization and finalization processes.

#### Configuration
Update the addr variable with the correct address of the Oakestra Root server.

#### Code Overview
The client establishes a gRPC connection to the Oakestra Root server and invokes the following methods:
HandleInitGreeting: Initiates the registration process by sending a greeting message.
HandleInitFinal: Completes the registration process by providing additional information, such as manager and network component ports, cluster name, cluster location, and custom cluster information.




## Hardware Aggregation

Thes directory [aggregation](./aggregation) provides an integration between a Kubernetes cluster and Oakestra Root. It collects information about the cluster's resource usage and sends it to an Oakestra Root server.

#### Functionality
The program aggregates information about the Kubernetes cluster, including CPU and memory usage, GPU information (if available), and general node details. This information is then formatted into JSON and sent to the specified Oakestra Root server.

#### Configuration
Ensure the correct path to the kubeconfig file is set using the --kubeconfig flag. This can be adjusted in the flag.StringVar function call.
Customize the Oakestra Root server address by updating the url variable in the SendClusterInfoToRoot function.

#### Getting started
1. Ensure you have Go installed on your machine.

2. Set up your kubectl configuration to point to the desired Kubernetes cluster.

3. Locate your kubeconfig file and ensure it's accessible.

4. Run the program using the following command:

    ``` bash
    go run .
    ```

## Oakestra Kubernetes Controller

The `kubenetesProxy` package provides functionalities for interacting with the Oakestra Kubernetes Controller.

### Configuration
Make sure to configure your Kubernetes client by calling the `SetUpKubernetesController` function before using other functionalities.

### Functions
- `CreateOakestraJob`: Creates a new OakestraJob based on the provided data, updating it if it already exists.
- `DeleteInstance`: Deletes an instance from the specified OakestraJob.

## Oakestra Job Controller