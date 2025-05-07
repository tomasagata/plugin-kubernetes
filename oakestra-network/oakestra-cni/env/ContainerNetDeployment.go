package env

import (
	"oakestra/cni-plugin/logger"
	// "oakestra/cni-plugin/mqtt"
	// "oakestra/cni-plugin/network"
	// "fmt"
	"log"
	// "net"
	"oakestra/cni-plugin/models"
	"path/filepath"
	// "runtime/debug"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

// type ContainerDeyplomentHandler struct {
// 	env *Environment
// }

// var containerHandler *ContainerDeyplomentHandler = nil

// func GetContainerNetDeployment() *ContainerDeyplomentHandler {
// 	if containerHandler == nil {
// 		logger.ErrorLogger().Fatal("Container Handler not initialized")
// 	}
// 	return containerHandler
// }

// func InitContainerDeployment(env *Environment) {
// 	containerHandler = &ContainerDeyplomentHandler{
// 		env: env,
// 	}
// }

// AttachNetworkToContainer Attach a Docker container to the bridge and the current network environment
func DeployNetwork(deployment *models.ConnectNetworkResponse) error {

	netnsPath := filepath.Join("/var/run/netns", deployment.ContainerNetNs)
	fd, err := unix.Open(netnsPath, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		log.Printf("File Descripter Error: %v", err)
	}

	// env := h.env
	cleanup := func(veth *netlink.Veth) {
		_ = netlink.LinkDel(veth)
	}

	vethIfce, err := createVethsPairAndAttachToBridge(deployment)
	if err != nil {
		go cleanup(vethIfce)
		return err
	}

	// Attach veth2 to the docker container
	logger.DebugLogger().Println("Attaching peerveth to container ")
	peerVeth, err := netlink.LinkByName(vethIfce.PeerName)
	if err != nil {
		cleanup(vethIfce)
		return err
	}

	// if pid == 0 {
		if err := netlink.LinkSetNsFd(peerVeth, fd); err != nil {
			cleanup(vethIfce)
			return err
		}
	// } else {
	// 	if err := netlink.LinkSetNsPid(peerVeth, pid); err != nil {
	// 		cleanup(vethIfce)
	// 		return err
	// 	}
	// }

	// //generate a new ip for this container
	// ip, err := env.generateAddress()
	// if err != nil {
	// 	cleanup(vethIfce)
	// 	return err
	// }

	// // generate a new ipv6 for this container
	// ipv6, err := env.generateIPv6Address()
	// if err != nil {
	// 	cleanup(vethIfce)
	// 	env.freeContainerAddress(ip)
	// 	return err
	// }

	// set ip to the container veth
	logger.DebugLogger().Println("Assigning ip ", deployment.ContainerIP+deployment.HostBridgeIPMask, " to container ")
	if err := addPeerLinkNetwork(netnsPath, deployment.ContainerIP+deployment.HostBridgeIPMask, vethIfce.PeerName); err != nil {
		logger.ErrorLogger().Println("Error in addPeerLinkNetwork")
		cleanup(vethIfce)
		// env.freeContainerAddress(ip)
		// env.freeContainerAddress(ipv6)
		return err
	}

	logger.DebugLogger().Println("Disabling DAD for IPv6")
	if err := disableDAD(netnsPath, vethIfce.PeerName); err != nil {
		logger.ErrorLogger().Println("Error in Disabling DAD")
		cleanup(vethIfce)
		// env.freeContainerAddress(ip)
		// env.freeContainerAddress(ipv6)
		return err
	}

	logger.DebugLogger().Println("Assigning ipv6 ", deployment.ContainerIPv6+deployment.HostBridgeIPv6Mask, " to container ")

	if err := addPeerLinkNetwork(netnsPath, deployment.ContainerIPv6+deployment.HostBridgeIPv6Mask, vethIfce.PeerName); err != nil {
		logger.ErrorLogger().Println("Error in addPeerLinkNetworkv6")
		cleanup(vethIfce)
		// env.freeContainerAddress(ip)
		// env.freeContainerAddress(ipv6)
		return err
	}

	// Add traffic route to bridge
	logger.DebugLogger().Println("Setting container routes ")
	if err = setContainerRoutes(netnsPath, vethIfce.PeerName, deployment.HostBridgeIP); err != nil {
		logger.ErrorLogger().Println("Error in setContainerRoutes")
		cleanup(vethIfce)
		// env.freeContainerAddress(ip)
		// env.freeContainerAddress(ipv6)
		return err
	}

	if err = setIPv6ContainerRoutes(netnsPath, vethIfce.PeerName, deployment.HostBridgeIPv6); err != nil {
		logger.ErrorLogger().Println("Error in setIPv6ContainerRoutes")
		cleanup(vethIfce)
		// env.freeContainerAddress(ip)
		// env.freeContainerAddress(ipv6)
		return err
	}

	// env.BookVethNumber()
	// if err = env.setVethFirewallRules(vethIfce.Name); err != nil {
	// 	logger.ErrorLogger().Println("Error in setFirewallRules")
	// 	cleanup(vethIfce)
	// 	env.freeContainerAddress(ip)
	// 	env.freeContainerAddress(ipv6)
	// 	return err
	// }

	// if err = network.ManageContainerPorts(ip, portmapping, network.OpenPorts); err != nil {
	// 	logger.ErrorLogger().Println("Error in ManageContainerPorts v4")
	// 	debug.PrintStack()
	// 	cleanup(vethIfce)
	// 	env.freeContainerAddress(ip)
	// 	env.freeContainerAddress(ipv6)
	// 	return err
	// }

	// if err = network.ManageContainerPorts(ipv6, portmapping, network.OpenPorts); err != nil {
	// 	logger.ErrorLogger().Println("Error in ManageContainerPorts v6")
	// 	debug.PrintStack()
	// 	cleanup(vethIfce)
	// 	env.freeContainerAddress(ip)
	// 	env.freeContainerAddress(ipv6)
	// 	return err
	// }

	// env.deployedServicesLock.Lock()
	// env.deployedServices[fmt.Sprintf("%s.%d", sname, instancenumber)] = service{
	// 	ip:          ip,
	// 	ipv6:        ipv6,
	// 	sname:       sname,
	// 	portmapping: portmapping,
	// 	veth:        vethIfce,
	// }
	// env.deployedServicesLock.Unlock()
	// logger.DebugLogger().Printf("New deployedServices table: %v", env.deployedServices)
	return nil
}

func DetachContainer(deployment *models.DettachNetworkResponse) error {
// func (env *Environment) DetachContainer(sname string, instance int) {
	// snameAndInstance := fmt.Sprintf("%s.%d", sname, instance)
	// env.deployedServicesLock.RLock()
	// s, ok := env.deployedServices[snameAndInstance]
	// env.deployedServicesLock.RUnlock()
	// if ok {
	// 	_ = env.translationTable.RemoveByNsip(s.ip)
	// 	env.deployedServicesLock.Lock()
	// 	delete(env.deployedServices, snameAndInstance)
	// 	env.deployedServicesLock.Unlock()
	// 	env.freeContainerAddress(s.ip)
	// 	env.freeContainerAddress(s.ipv6)
	// 	_ = network.ManageContainerPorts(s.ip, s.portmapping, network.ClosePorts)
	// 	_ = network.ManageContainerPorts(s.ipv6, s.portmapping, network.ClosePorts)
	veth, err := netlink.LinkByName(deployment.VethPeer1Name)
	if err != nil {
		logger.ErrorLogger().Println("DetachContainer: Error in LinkByName")
		return err
	}
	err = netlink.LinkDel(veth)
	if err != nil {
		logger.ErrorLogger().Println("DetachContainer: Error in LinkDel")
		return err
	}
	// 	// if no interest registered delete all remaining info about the service
	// 	if !mqtt.MqttIsInterestRegistered(sname) {
	// 		env.RemoveServiceEntries(sname)
	// 	}
	// }
	return nil
}
