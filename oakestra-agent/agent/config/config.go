package config

import (
	"log"
	"os"
	"strconv"
)

func GetConfig() *Config {
	config := &Config{ // Default values
		RootSystemManagerPort:   10000,
		RootServiceManagerPort:  10099,
		RootGRPCPort:            50052,
		NetworkComponentPort:    10110,
		MyPort:                  10100,
		NodePort:                30000,
		ClusterServiceManagerIP: "localhost",
		ClusterName:             "k8s",
		ClusterLocation:         "Munich",
	}
	config.readENVVariables()

	return config
}

type Config struct {
	RootSystemManagerIP     string
	ClusterName             string
	ClusterLocation         string
	ClusterServiceManagerIP string
	ClusterID               string

	NodePort               int
	MyPort                 int
	NetworkComponentPort   int
	RootGRPCPort           int
	RootSystemManagerPort  int
	RootServiceManagerPort int
}

func (c *Config) readENVVariables() {

	if temp := os.Getenv("ROOT_SYSTEM_MANAGER_IP"); temp != "" {
		c.RootSystemManagerIP = temp
	} else {
		log.Fatal("ERROR: ROOT_SYSTEM_MANAGER_IP environment variable is not set")
	}

	if temp := os.Getenv("ROOT_SYSTEM_MANAGER_PORT"); temp != "" {
		if RootSystemManagerPort, err := convertStringToInt(temp); err != nil {
			log.Fatal("ERROR: ROOT_SYSTEM_MANAGER_PORT environment variable not in right format")
		} else {
			c.RootSystemManagerPort = RootSystemManagerPort
		}
	} else {
		logVariableWarning("ROOT_SYSTEM_MANAGER_PORT", "10000")
	}

	if temp := os.Getenv("CLUSTER_SERVICE_MANAGER_IP"); temp != "" {
		c.ClusterServiceManagerIP = temp
	} else {
		logVariableWarning("CLUSTER_SERVICE_MANAGER_IP", "localhost")
	}

	if temp := os.Getenv("ROOT_GRPC_PORT"); temp != "" {
		if rootGRPCPort, err := convertStringToInt(temp); err != nil {
			log.Fatal("ERROR: ROOT_GRPC_PORT environment variable not in right format")
		} else {
			c.RootGRPCPort = rootGRPCPort
		}
	} else {
		logVariableWarning("ROOT_GRPC_PORT", "50052")
	}

	if temp := os.Getenv("NODE_PORT"); temp != "" {
		if NodePort, err := convertStringToInt(temp); err != nil {
			log.Fatal("ERROR: NODE_PORT environment variable not in right format")
		} else {
			c.NodePort = NodePort
		}
	} else {
		logVariableWarning("NODE_PORT", "30000")
	}

	if temp := os.Getenv("CLUSTER_SERVICE_MANAGER_PORT"); temp != "" {
		if NetworkComponentPort, err := convertStringToInt(temp); err != nil {
			log.Fatal("ERROR: CLUSTER_SERVICE_MANAGER_PORT environment variable not in right format")
		} else {
			c.NetworkComponentPort = NetworkComponentPort
		}
	} else {
		logVariableWarning("CLUSTER_SERVICE_MANAGER_PORT", "10110")
	}

	if temp := os.Getenv("ROOT_SERIVCE_MANAGER_PORT"); temp != "" {
		if RootServiceManagerPort, err := convertStringToInt(temp); err != nil {
			log.Fatal("ERROR: ROOT_SERIVCE_MANAGER_PORT environment variable not in right format")
		} else {
			c.RootServiceManagerPort = RootServiceManagerPort
		}
	} else {
		logVariableWarning("ROOT_SERIVCE_MANAGER_PORT", "10099")
	}

	if temp := os.Getenv("MY_PORT"); temp != "" {
		if MyPort, err := convertStringToInt(temp); err != nil {
			log.Fatal("ERROR: MY_PORT environment variable not in right format")
		} else {
			c.MyPort = MyPort
		}
	} else {
		logVariableWarning("MY_PORT", "10100")
	}

	if temp := os.Getenv("CLUSTER_LOCATION"); temp != "" {
		c.ClusterLocation = temp
	} else {
		logVariableWarning("CLUSTER_LOCATION ", "Munich")
	}

	if temp := os.Getenv("CLUSTER_NAME"); temp != "" {
		c.ClusterName = temp
	} else {
		logVariableWarning("CLUSTER_NAME", "k8s")
	}
}

func (c *Config) SetClusterID(clusterID string) {
	c.ClusterID = clusterID
}

func logVariableWarning(variableName, defaultValue string) {
	log.Printf("WARNING: %s environment variable is not set or could not be parsed, using default value %s\n",
		variableName, defaultValue)
}

func convertStringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}
