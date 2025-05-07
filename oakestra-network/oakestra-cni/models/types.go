package models

type ConnectNetworkRequest struct {
	NetworkNamespace string `json:"networkNamespace"`
	Servicename      string `json:"servicename"`
	Instancenumber   int    `json:"instancenumber"`
	PodName          string `json:"podName"`
}

type DettachNetworkRequest struct {
	Servicename    string `json:"serviceName"`
	Instancenumber int    `json:"instanceNumber"`
}

type ConnectNetworkResponse struct {
	ServiceName string `json:"serviceName"`
	HostVethName string `json:"hostVethName"`
	HostBridgeName string `json:"hostBridgeName"`
	HostBridgeIP string `json:"hostBridgeIP"`
	HostBridgeIPMask string `json:"hostBridgeIPMask"`
	HostBridgeIPv6 string `json:"hostBridgeIPv6"`
	HostBridgeIPv6Mask string `json:"hostBridgeIPv6Mask"`
	ContainerVethName string `json:"containerVethName"`
	ContainerNetNs string `json:"containerNetNs"`
	ContainerIP string `json:"containerIP"`
	ContainerIPv6 string `json:"containerIPv6"`
	Mtu       int    `json:"mtu"`
}

type DettachNetworkResponse struct {
	VethPeer1Name string `json:"vethPeer1Name"`
}