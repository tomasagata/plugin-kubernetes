package hardwareagreggation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	config "oakestra/plugin-kubernetes/oakestra-agent/agent/config"
	k8sclient "oakestra/plugin-kubernetes/oakestra-agent/agent/kubernetesClient"
	"os"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type ClusterInfo struct {
	CPUPercent                 int      `json:"cpu_percent"`
	CPUCores                   int      `json:"cpu_cores"`
	CumulativeMemoryInMB       int      `json:"cumulative_memory_in_mb"`
	MemoryPercent              int      `json:"memory_percent"`
	NumberOfNodes              int      `json:"number_of_nodes"`
	Virtualization             []string `json:"virtualization"`
	More                       int      `json:"more"`
	Jobs                       []Job    `json:"jobs"`
	WorkerGroups               string   `json:"worker_groups"`
	AggregationPerArchitecture []Arch   `json:"aggregation_per_architecture"`

	// Not used, no GPU support so far.
	GPUCores   int      `json:"gpu_cores"`
	GPUPercent int      `json:"gpu_percent"`
	GPUDrivers []string `json:"gpu_drivers"`
	GPUTemp    int      `json:"gpu_temp"`
	GPUTotMem  int      `json:"gpu_tot_mem"`
	GPUMemUsed int      `json:"gpu_mem_used"`
}

// More architectures can be added in the future.
type Arch struct {
	AMD AggArchitecture `json:"amd"`
}
type AggArchitecture struct {
	CPUPercent int `json:"cpu_percent"`
	CPUCores   int `json:"cpu_cores"`
	Memory     int `json:"memory"`
	MemoryInMB int `json:"memory_in_mb"`
}

type Job struct {
	ID           string     `json:"_id"`
	SystemJobID  string     `json:"system_job_id"`
	JobName      string     `json:"job_name"`
	Status       string     `json:"status"`
	InstanceList []Instance `json:"instance_list"`
}

type Instance struct {
	InstanceNumber int64  `json:"instance_number"`
	Status         string `json:"status"`
	HostPort       string `json:"hostPort"`
	HostIP         string `json:"hostIP"`
	WorkerID       string `json:"workerID"`
}

var (
	clientSet     *kubernetes.Clientset
	metricsClient *metricsv.MetricsV1beta1Client
)

func getJobsOfCluster(k8sClient k8sclient.KubernetesClient) ([]Job, error) {
	jobs := []Job{}

	oakestraJobs, err := k8sClient.GetExistingOakestraJobs()
	if err != nil {
		return nil, fmt.Errorf("getJobsOfCluster Error: %v", err)
	}

	for _, oakestraObj := range oakestraJobs.Items {
		instanceList := []Instance{}

		labels := oakestraObj.GetLabels()
		id := labels["ID"]

		obj := oakestraObj.UnstructuredContent()
		uid := obj["metadata"].(map[string]interface{})["uid"].(string)
		name := obj["metadata"].(map[string]interface{})["name"].(string)

		spec, parseOk := oakestraObj.Object["spec"].(map[string]interface{})
		if !parseOk {
			return nil, errors.New("error in parsing OakestraJob spec")
		}
		instanceListSpec, parseOk := spec["instance_list"].([]interface{})
		if !parseOk {
			return nil, errors.New("error in parsing OakestraJob spec: instance_list not found or invalid type")
		}

		for _, instance := range instanceListSpec {
			instanceMap, ok := instance.(map[string]interface{})
			if !ok {
				return nil, errors.New("error in parsing instance_list: invalid instance type")
			}
			// Test, if attributes have a value.
			if val := instanceMap["host_IP"]; val == nil {
				continue
			}
			instanceList = append(instanceList, Instance{
				HostIP:         instanceMap["host_IP"].(string),
				HostPort:       instanceMap["host_port"].(string),
				InstanceNumber: instanceMap["instance_number"].(int64),
				Status:         instanceMap["status"].(string),
				WorkerID:       instanceMap["worker_ID"].(string),
			})
		}

		state, parseOk := spec["state"].(string)
		if !parseOk {
			state = ""
			log.Print("Warning: error in parsing OakestraJob spec: State not found or not a string")
		}
		jobs = append(jobs, Job{ID: uid, SystemJobID: id, JobName: name, Status: state, InstanceList: instanceList})
	}

	return jobs, nil
}

func SetUp() error {
	var config *rest.Config
	var err error

	if _, inCluster := os.LookupEnv("KUBERNETES_SERVICE_HOST"); inCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Error building in-cluster kubeconfig: %v\n", err)
		}
	}
	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {

		return fmt.Errorf("create k8s client failed: %w", err)
	}

	metricsClient, err = metricsv.NewForConfig(config)
	if err != nil {

		return fmt.Errorf("create metricsclient failed: %w", err)
	}

	log.Println("Hardware Aggregation successfully set up")

	return nil
}

func StartBackgroundServiceClusterInformation(timeInterval time.Duration,
	config *config.Config, k8sClient k8sclient.KubernetesClient) {
	log.Println("Start Hardware Aggregation Updates to Root")
	rootSystemManagerAddress := config.RootSystemManagerIP + ":" + strconv.Itoa(config.RootSystemManagerPort)

	go func() {
		for {
			clusterInfo, err := aggregateNodeInfo(k8sClient)
			if err != nil {
				log.Fatalf("Error in aggregating Node Info: %s", err)
			}

			err = sendClusterInfoToRoot(clusterInfo, rootSystemManagerAddress, config.ClusterID)
			if err != nil {
				log.Printf("Could not send hardware information to root: %s", err)
				log.Printf("Retrying...")
			}

			time.Sleep(timeInterval)
		}
	}()
}

func aggregateNodeInfo(k8sClient k8sclient.KubernetesClient) (ClusterInfo, error) {
	cumulativeCPU := 0
	cumulativeCPUCores := 0
	cumulativeMemory := 0
	cumulativeMemoryInMb := 0
	numberOfActiveNodes := 0
	technology := []string{"docker"}

	nodes, err := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, node := range nodes.Items {
		nodeMetrics, err := metricsClient.NodeMetricses().Get(context.TODO(), node.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}

		freeMemoryInMb := (node.Status.Capacity.Memory().Value() /
			(1024 * 1024)) - (nodeMetrics.Usage.Memory().Value() / (1024 * 1024))
		percentageUsageCPU := float32(nodeMetrics.Usage.Cpu().MilliValue()) / float32(node.Status.Capacity.Cpu().MilliValue())
		freeCPUCore := (node.Status.Capacity.Cpu().MilliValue() - nodeMetrics.Usage.Cpu().MilliValue()) / 1000
		percentageUsageMemory := float32(nodeMetrics.Usage.Memory().MilliValue()) /
			float32(node.Status.Capacity.Memory().MilliValue())

		// GPU

		// gpuResources, exists := node.Status.Capacity["nvidia.com/gpu"]
		// if !exists {
		// 	log.Printf("No GPU resources found for node %s\n", node.Name)
		// } else {
		// 	log.Printf("GPU Count: %s\n", gpuResources.String())
		// 	// Add gpu_cores, gpu_info, and gpu_percent
		// }

		cumulativeMemoryInMb += int(freeMemoryInMb)
		cumulativeCPU += int(percentageUsageCPU * 100)
		cumulativeCPUCores += int(freeCPUCore)
		cumulativeMemory += int(percentageUsageMemory * 100)
		numberOfActiveNodes++
	}

	jobs, err := getJobsOfCluster(k8sClient)
	if err != nil {
		return ClusterInfo{}, err
	}

	return ClusterInfo{
		CPUPercent:           cumulativeCPU,
		CPUCores:             cumulativeCPUCores,
		CumulativeMemoryInMB: cumulativeMemoryInMb,
		MemoryPercent:        cumulativeMemory,
		NumberOfNodes:        2,
		Virtualization:       technology,
		More:                 0,
		Jobs:                 jobs,

		// So far only amd is supported.
		AggregationPerArchitecture: []Arch{
			{AMD: AggArchitecture{
				CPUPercent: cumulativeCPU,
				CPUCores:   cumulativeCPUCores,
				Memory:     cumulativeMemory,
				MemoryInMB: cumulativeMemoryInMb},
			},
		},

		// Default values to 0
		GPUCores:   0,
		GPUPercent: 0,
		GPUDrivers: []string{},
		GPUTemp:    0,
		GPUTotMem:  0,
		GPUMemUsed: 0,
	}, nil
}

func sendClusterInfoToRoot(clusterInfo ClusterInfo, rootSystemManagerAddress string, clusterID string) error {
	jsonData, err := json.Marshal(clusterInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal clusterInfo to JSON: %w", err)
	}

	u := &url.URL{
		Scheme: "http",
		Host:   rootSystemManagerAddress,
		Path:   fmt.Sprintf("/api/information/%s", clusterID),
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var ErrNonOKStatus = fmt.Errorf("server returned non-OK status")

		return fmt.Errorf("%w: %d", ErrNonOKStatus, resp.StatusCode)
	}

	return nil
}
