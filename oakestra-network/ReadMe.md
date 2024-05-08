# Oakestra Network

Dieses ReadMe erklärt die Funktionsweise und das Deployment des Oakestra Networks, das zusätzlich zum Default Netzwerk von Kubernetes gestartet wird. 


1. [Deployment](#oakestra-agent)
2. [Functionality](#Functionality)



## Deployment
Um das Oakestra Network zu starten, müssen mehrere Komponenten gestartet werden. Im ReadMe des Repos stehen bereits die Befehle. Hier werden nur die Env Variablen erklärt. 

Für den [Nodenetmanager](Deployment/oakestra-nodenetmanager/node-netmanager.yaml) müssen diese Variablen gesetzt werden. 


| Variable Name                   | Default Values      | Description                                                     | 
|---------------------------------|---------------------|-----------------------------------------------------------------|
| NODE_PORT          | 50103     | public node port for oakestra network                                                |
| MOSQUITTO_SVC_SERVICE_PORT       | 30033               | -                                             |
| MOSQUITTO_SVC_SERVICE_HOST       | *Needs to be set*               | NodePort of one kubernetes node


Für den [Cluster Service Manager](Deployment/oakestra-cluster-service-manager/oakestra-cluster-service-manager.yaml) müssen diese Variablen gesetzt werden. 

The **MQTT_BROKER_URL** and **CLUSTER_MONGO_URL** need to be configured. These values should be derived from the **mosquitto-svc** and **mongo-svc** services located in the *oakestra-system* namespace. The assigned values must match the ClusterIPs of these services.


| Variable Name               | Default Value          | Description  |
|-----------------------------|------------------------|--------------|
| MY_PORT                     | 10110                  | Local port which starts server  |
| MQTT_BROKER_PORT            | 10003                  |       -       |
| MQTT_BROKER_URL             | *Needs to be set*        |  ClusterIP of Mosquitto Service     |
| ROOT_SERVICE_MANAGER_URL    | *Needs to be set*       |    IP Oakestra Root Network          |
| ROOT_SERVICE_MANAGER_PORT   | 10099                  |   Port Oakestra Root Network            |
| SYSTEM_MANAGER_URL          | *Needs to be set*       |     IP Oakestra Root         |
| SYSTEM_MANAGER_PORT         | 10000                  |    Port Oakestra Root          |
| CLUSTER_MONGO_URL           | *Needs to be set*         |   ClusterIP of MongoDB Service    |
| CLUSTER_MONGO_PORT          | 27017                  |        -      |




## Functionality

The Oakestra network components for Kubernetes are elaborated in greater detail in the Oakestra wiki. The system is composed of three primary components.


### Oakestra CNI 
A dedicated CNI is required for the Oakestra network, which, in conjunction with the netmanager, adjusts the Network Namespace of the containers and facilitates connectivity to other Oakestra resources.

### Node Netmanager
This manager operates on each node and functions as a proxy to facilitate connections to other Oakestra networks.

### Oakestra Cluster Service Manager
In this component, one can easily deploy the default Oakestra image. Additionally, it is necessary to launch both Mosquitto and MongoDB.


