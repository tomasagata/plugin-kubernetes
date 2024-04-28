// Package main implements a client for Greeter service.
package kubenetesproxy

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	pb "oakestra/plugin-kubernetes/oakestra-agent/agent/clusterRegistration"
	config "oakestra/plugin-kubernetes/oakestra-agent/agent/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Register(config *config.Config) (string, error) {
	rootSystemManagerAddress := config.RootSystemManagerIP + ":" + strconv.Itoa(config.RootGRPCPort)

	fmt.Println(rootSystemManagerAddress)
	conn, err := grpc.Dial(rootSystemManagerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", fmt.Errorf("GRPC Dial failed: %w", err)
	}
	defer conn.Close()
	c := pb.NewRegisterClusterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.HandleInitGreeting(ctx, &pb.CS1Message{HelloServiceManager: "This is a K8s Cluster"})
	if err != nil {
		return "", fmt.Errorf("HandleInitGreeting failed: %w", err)
	}
	log.Printf("Answer of HandleInitGreeting: %s", r.GetHelloClusterManager())

	cs2Message := &pb.CS2Message{
		ManagerPort:          int32(config.NodePort), // NodePort for Kubernetes Pod
		NetworkComponentPort: int32(config.NetworkComponentPort),
		ClusterName:          config.ClusterName,
		ClusterLocation:      config.ClusterLocation,
	}

	clusterInfo := []*pb.KeyValue{
		{Key: "cluster_type", Value: "k8s"},
	}

	cs2Message.ClusterInfo = clusterInfo

	r2, err := c.HandleInitFinal(ctx, cs2Message)
	if err != nil {
		return "", fmt.Errorf("HandleInitFinal failed: %w", err)
	}

	clusterID := r2.GetId()

	log.Printf("Received ClusterID: %s", clusterID)

	return clusterID, nil
}
