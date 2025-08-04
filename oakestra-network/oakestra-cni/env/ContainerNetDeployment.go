package env

import (
	"oakestra/cni-plugin/logger"
	"log"
	"oakestra/cni-plugin/models"
	"path/filepath"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)


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

	if err := netlink.LinkSetNsFd(peerVeth, fd); err != nil {
		cleanup(vethIfce)
		return err
	}

	// set ip to the container veth
	logger.DebugLogger().Println("Assigning ip ", deployment.ContainerIP+deployment.HostBridgeIPMask, " to container ")
	if err := addPeerLinkNetwork(netnsPath, deployment.ContainerIP+deployment.HostBridgeIPMask, vethIfce.PeerName); err != nil {
		logger.ErrorLogger().Println("Error in addPeerLinkNetwork")
		cleanup(vethIfce)
		return err
	}

	logger.DebugLogger().Println("Disabling DAD for IPv6")
	if err := disableDAD(netnsPath, vethIfce.PeerName); err != nil {
		logger.ErrorLogger().Println("Error in Disabling DAD")
		cleanup(vethIfce)
		return err
	}

	logger.DebugLogger().Println("Assigning ipv6 ", deployment.ContainerIPv6+deployment.HostBridgeIPv6Mask, " to container ")

	if err := addPeerLinkNetwork(netnsPath, deployment.ContainerIPv6+deployment.HostBridgeIPv6Mask, vethIfce.PeerName); err != nil {
		logger.ErrorLogger().Println("Error in addPeerLinkNetworkv6")
		cleanup(vethIfce)
		return err
	}

	// Add traffic route to bridge
	logger.DebugLogger().Println("Setting container routes ")
	if err = setContainerRoutes(netnsPath, vethIfce.PeerName, deployment.HostBridgeIP); err != nil {
		logger.ErrorLogger().Println("Error in setContainerRoutes")
		cleanup(vethIfce)
		return err
	}

	if err = setIPv6ContainerRoutes(netnsPath, vethIfce.PeerName, deployment.HostBridgeIPv6); err != nil {
		logger.ErrorLogger().Println("Error in setIPv6ContainerRoutes")
		cleanup(vethIfce)
		return err
	}

	return nil
}

func DetachContainer(deployment *models.DettachNetworkResponse) error {
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
	return nil
}
