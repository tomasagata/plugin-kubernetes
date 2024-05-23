package kubenetesclient

import (
	"context"
	"errors"
	"fmt"
	"log"
	"oakestra/plugin-kubernetes/oakestra-agent/agent/config"
	"os"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type OakestraJSONRequest struct {
	AddedFiles                    []string              `json:"added_files"`
	ApplicationName               string                `json:"app_name"`
	ApplicationNamespace          string                `json:"app_ns"`
	ApplicationID                 string                `json:"applicationID"`
	BandwidthIn                   int                   `json:"bandwidth_in"`
	BandwidthOut                  int                   `json:"bandwidth_out"`
	Cmd                           []string              `json:"cmd"`
	Code                          string                `json:"code"`
	Image                         string                `json:"image"`
	InstanceList                  []InstanceJSONRequest `json:"instance_list"`
	JobName                       string                `json:"job_name"`
	Memory                        int                   `json:"memory"`
	MicroserviceID                string                `json:"microserviceID"`
	MicroserviceName              string                `json:"microservice_name"`
	MicroserviceNamespace         string                `json:"microservice_namespace"`
	NextInstanceProgressiveNumber int16                 `json:"next_instance_progressive_number"`
	Port                          string                `json:"port"`
	State                         string                `json:"state"`
	Status                        string                `json:"status"`
	StatusDetail                  string                `json:"status_detail"`
	Storage                       int                   `json:"storage"`
	VCPUs                         int                   `json:"vcpus"`
	VGPUs                         int                   `json:"vgpus"`
	Disk                          int                   `json:"disk"`
	Virtualization                string                `json:"virtualization"`
	VTPUs                         int                   `json:"vtpus"`
	Environment                   []string              `json:"environment"`
}

type InstanceJSONRequest struct {
	ClusterID       string `json:"cluster_id"`
	ClusterLocation string `json:"cluster_location"`
	InstanceNumber  int    `json:"instance_number"`
}

type OakestraJobSpec struct {
	JobName                       string                `json:"job_name"`
	AddedFiles                    []string              `json:"added_files"`
	ApplicationName               string                `json:"app_name"`
	ApplicationNamespace          string                `json:"app_ns"`
	ApplicationID                 string                `json:"applicationID"`
	BandwidthIn                   int                   `json:"bandwidth_in"`
	BandwidthOut                  int                   `json:"bandwidth_out"`
	Cmd                           []string              `json:"cmd"`
	Code                          string                `json:"code"`
	Image                         string                `json:"image"`
	InstanceList                  []InstanceJSONRequest `json:"instance_list"`
	Memory                        int                   `json:"memory"`
	MicroserviceID                string                `json:"microserviceID"`
	MicroserviceName              string                `json:"microservice_name"`
	MicroserviceNamespace         string                `json:"microservice_namespace"`
	NextInstanceProgressiveNumber int16                 `json:"next_instance_progressive_number"`
	Port                          string                `json:"port"`
	State                         string                `json:"state"`
	Storage                       int                   `json:"storage"`
	VCPUs                         int                   `json:"vcpus"`
	VGPUs                         int                   `json:"vgpus"`
	Virtualization                string                `json:"virtualization"`
	VTPUs                         int                   `json:"vtpus"`
	Status                        string                `json:"status"`
	StatusDetail                  string                `json:"status_detail"`
	Disk                          int                   `json:"disk"`

	TargetNode           string   `json:"target_node"`
	VMImages             []string `json:"vm_images"`
	Arch                 []string `json:"arch"`
	Args                 []string `json:"args"`
	Environment          []string `json:"environment"`
	SLAViolationStrategy string   `json:"sla_violation_strategy"`
}

type OakestraJob struct {
	APIVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Metadata   Metadata        `json:"metadata"`
	Spec       OakestraJobSpec `json:"spec"`
}

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type KubernetesClient struct {
	dynamicClient     *dynamic.DynamicClient
	clientSet         *kubernetes.Clientset
	resourceNamespace string
}

func NewKubernetesClient() (KubernetesClient, error) {
	var config *rest.Config
	var err error

	config, err = rest.InClusterConfig()
	if err != nil {
		fmt.Printf("Error building in-cluster kubeconfig: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Running in a Kubernetes cluster.")

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return KubernetesClient{}, fmt.Errorf("configure kubernetes client failed: %w", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return KubernetesClient{}, fmt.Errorf("configure dynamicClient failed: %w", err)
	}

	return KubernetesClient{
		dynamicClient:     dynamicClient,
		clientSet:         clientSet,
		resourceNamespace: "oakestra",
	}, nil
}

func (c *KubernetesClient) UpdateClusterInfoConfigMap(clusterID string, config *config.Config) error {
	configMapName := "oakestra-cluster-info"
	namespace := "oakestra-controller-manager"
	keyClusterID := "CLUSTER_ID"
	keyRootIP := "OAKESTRA_ROOT_IP"
	keyRootNetworkPort := "OAKESTRA_ROOT_NETWORK_PORT"

	configMap, err := c.clientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get ConfigMap '%s': %w", configMapName, err)
	}

	configMap.Data[keyClusterID] = clusterID
	configMap.Data[keyRootIP] = config.RootSystemManagerIP
	configMap.Data[keyRootNetworkPort] = strconv.Itoa(config.RootServiceManagerPort)

	_, err = c.clientSet.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update ConfigMap '%s': %w", configMapName, err)
	}

	return nil
}

func (c KubernetesClient) CreateOakestraJob(jsonRequest OakestraJSONRequest, systemJobID string,
	instanceNumber int) (OakestraJob, error) {
	oakestraJob := createOakestraJobStructFromBody(jsonRequest)

	oakestraObj, err := c.getExistingOakestraJob(systemJobID)

	if err != nil {
		err = c.createOakestraJobInstance(oakestraJob)
		if err != nil {
			return OakestraJob{}, fmt.Errorf("failed to deploy OakestraJob: %w", err)
		}
		log.Printf("Custom Resource created successfully")
	} else {
		err = c.updateOakestraJobInstance(oakestraObj, oakestraJob, instanceNumber)
		if err != nil {
			return OakestraJob{}, fmt.Errorf("failed to update existing OakestraJob: %w", err)
		}
		log.Printf("Custom Resource updated successfully")
	}

	return oakestraJob, nil
}

func createOakestraJobStructFromBody(jsonRequest OakestraJSONRequest) OakestraJob {
	oakestraJob := OakestraJob{
		APIVersion: "oakestra.oakestra.kubernetes/v1",
		Kind:       "OakestraJob",
		Metadata: Metadata{
			Name:      jsonRequest.MicroserviceName,
			Namespace: "oakestra",
		},
		Spec: OakestraJobSpec{
			JobName:                       jsonRequest.JobName,
			ApplicationName:               jsonRequest.ApplicationName,
			ApplicationNamespace:          jsonRequest.ApplicationNamespace,
			ApplicationID:                 jsonRequest.ApplicationID,
			AddedFiles:                    jsonRequest.AddedFiles,
			BandwidthIn:                   jsonRequest.BandwidthIn,
			BandwidthOut:                  jsonRequest.BandwidthOut,
			Cmd:                           jsonRequest.Cmd,
			Image:                         jsonRequest.Image,
			Memory:                        jsonRequest.Memory,
			Port:                          jsonRequest.Port,
			Code:                          jsonRequest.Code,
			MicroserviceID:                jsonRequest.MicroserviceID,
			MicroserviceName:              jsonRequest.MicroserviceName,
			MicroserviceNamespace:         jsonRequest.MicroserviceNamespace,
			NextInstanceProgressiveNumber: jsonRequest.NextInstanceProgressiveNumber,
			State:                         jsonRequest.State,
			Status:                        "KUBERNETES_SCHEDULED",                         //jsonRequest.Status,
			StatusDetail:                  "Kubernetes takes care of managing the status", //,jsonRequest.StatusDetail,
			Storage:                       jsonRequest.Storage,
			VCPUs:                         jsonRequest.VCPUs,
			VGPUs:                         jsonRequest.VGPUs,
			Virtualization:                jsonRequest.Virtualization,
			VTPUs:                         jsonRequest.VTPUs,
			Disk:                          jsonRequest.Disk,
			InstanceList:                  jsonRequest.InstanceList,
			Environment:                   jsonRequest.Environment,
		},
	}

	return oakestraJob
}

func (c KubernetesClient) GetExistingOakestraJobs() (*unstructured.UnstructuredList, error) {
	resourceList, err := c.dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "oakestra.oakestra.kubernetes",
		Version:  "v1",
		Resource: "oakestrajobs",
	}).Namespace(c.resourceNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve existing OakestraJobs: %v", err)
	}

	return resourceList, nil
}

func (c KubernetesClient) getExistingOakestraJob(id string) (*unstructured.Unstructured, error) {
	var item *unstructured.Unstructured

	options := metav1.ListOptions{
		LabelSelector: c.createLabelSelector("ID", id),
	}

	resourceList, err := c.dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "oakestra.oakestra.kubernetes",
		Version:  "v1",
		Resource: "oakestrajobs",
	}).Namespace(c.resourceNamespace).List(context.TODO(), options)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve existing OakestraJob: %v", err)
	}

	if len(resourceList.Items) > 0 {
		return &resourceList.Items[0], nil
	}

	return item, errors.New("OakestraJob doest not exist. Controller creates a new one")
}

func (c KubernetesClient) createOakestraJobInstance(oakestraJob OakestraJob) error {
	obj := convertOakestraJobToUnstructured(oakestraJob)

	result, err := c.dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "oakestra.oakestra.kubernetes",
		Version:  "v1",
		Resource: "oakestrajobs",
	}).Namespace("oakestra").Create(context.TODO(), obj, metav1.CreateOptions{})

	log.Printf("Return of Kubernetes Create: %v", result)

	if err != nil {
		return fmt.Errorf("create oakestraJob failed: %w", err)
	}
	return nil
}

func (c KubernetesClient) updateOakestraJobInstance(oakestraObj *unstructured.Unstructured,
	oakestraJob OakestraJob, instanceNumber int) error {
	spec, parseOk := oakestraObj.Object["spec"].(map[string]interface{})
	if !parseOk {

		return errors.New("error in parsing OakestraJob spec")
	}

	instanceList, parseOk := spec["instance_list"].([]interface{})
	if !parseOk {

		return errors.New("error in parsing OakestraJob spec: instance_list not found or invalid type")
	}

	newInstance := make(map[string]interface{})
	if len(instanceList) == 0 {
		newInstance["cluster_location"] = oakestraJob.Spec.InstanceList[0].ClusterLocation
		newInstance["cluster_ID"] = oakestraJob.Spec.InstanceList[0].ClusterID
	} else {
		// last item index
		lastItem, parseOk := instanceList[len(instanceList)-1].(map[string]interface{})
		if !parseOk {
			return errors.New("error in parsing OakestraJob spec: last item in instance_list is not a map")
		}
		for key, value := range lastItem {
			newInstance[key] = value
		}
	}

	newInstance["instance_number"] = instanceNumber

	instanceList = append(instanceList, newInstance)
	spec["instance_list"] = instanceList

	// Updating the OakestraJob resource
	_, err := c.dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "oakestra.oakestra.kubernetes",
		Version:  "v1",
		Resource: "oakestrajobs",
	}).Namespace(c.resourceNamespace).Update(context.TODO(), oakestraObj, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update OakestraJob resource: %w", err)
	}

	return nil
}

func convertOakestraJobToUnstructured(oakestraJob OakestraJob) *unstructured.Unstructured {
	unstructuredObj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "oakestra.oakestra.kubernetes/v1",
			"kind":       "OakestraJob",
			"metadata": map[string]interface{}{
				"name":      oakestraJob.Spec.MicroserviceName,
				"namespace": oakestraJob.Metadata.Namespace,
				"labels": map[string]interface{}{
					"ID": oakestraJob.Spec.MicroserviceID,
				},
			},
			"spec": map[string]interface{}{
				"job_name":                         oakestraJob.Spec.JobName,
				"added_files":                      oakestraJob.Spec.AddedFiles,
				"application_ID":                   oakestraJob.Spec.ApplicationID,
				"application_name":                 oakestraJob.Spec.ApplicationName,
				"application_namespace":            oakestraJob.Spec.ApplicationNamespace,
				"bandwidth_in":                     oakestraJob.Spec.BandwidthIn,
				"bandwidth_out":                    oakestraJob.Spec.BandwidthOut,
				"cmd":                              oakestraJob.Spec.Cmd,
				"code":                             oakestraJob.Spec.Code,
				"image":                            oakestraJob.Spec.Image,
				"memory":                           oakestraJob.Spec.Memory,
				"microservice_ID":                  oakestraJob.Spec.MicroserviceID,
				"microservice_name":                oakestraJob.Spec.MicroserviceName,
				"microservice_namespace":           oakestraJob.Spec.MicroserviceNamespace,
				"next_instance_progressive_number": oakestraJob.Spec.NextInstanceProgressiveNumber,
				"port":                             oakestraJob.Spec.Port,
				"state":                            oakestraJob.Spec.State,
				"storage":                          oakestraJob.Spec.Storage,
				"vcpus":                            oakestraJob.Spec.VCPUs,
				"vgpus":                            oakestraJob.Spec.VGPUs,
				"virtualization":                   oakestraJob.Spec.Virtualization,
				"vtpus":                            oakestraJob.Spec.VTPUs,
				"status":                           oakestraJob.Spec.Status,
				"status_detail":                    oakestraJob.Spec.StatusDetail,
				"disk":                             oakestraJob.Spec.Disk,
				"instance_list":                    convertInstanceDetailsToUnstructured(oakestraJob.Spec.InstanceList),
				"environment":                      oakestraJob.Spec.Environment,
			},
		},
	}

	return unstructuredObj
}

func convertInstanceDetailsToUnstructured(instanceDetails []InstanceJSONRequest) []interface{} {
	unstructuredList := make([]interface{}, 0, len(instanceDetails))
	for _, instance := range instanceDetails {
		instanceObj := map[string]interface{}{
			"instance_number":  instance.InstanceNumber,
			"cluster_location": instance.ClusterLocation,
			"cluster_ID":       instance.ClusterID,
		}
		unstructuredList = append(unstructuredList, instanceObj)
	}

	return unstructuredList
}

// Expand to multiple labels?
func (c KubernetesClient) createLabelSelector(label string, key string) string {
	return fmt.Sprintf("%s=%s", label, key)
}

func (c KubernetesClient) DeleteInstance(systemJobID string, instanceNumber int64) error {
	options := metav1.ListOptions{
		LabelSelector: c.createLabelSelector("ID", systemJobID),
	}

	resourceList, err := c.dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "oakestra.oakestra.kubernetes",
		Version:  "v1",
		Resource: "oakestrajobs",
	}).Namespace(c.resourceNamespace).List(context.TODO(), options)
	if err != nil {
		return fmt.Errorf("get oakestraJobs list failed: %w", err)
	}

	for _, item := range resourceList.Items {
		spec, ok := item.Object["spec"].(map[string]interface{})
		if !ok {
			continue
		}
		instanceList, ok := spec["instance_list"].([]interface{})
		if !ok {
			continue
		}

		var newInstanceList []interface{}

		for i, inst := range instanceList {
			instance, ok := inst.(map[string]interface{})
			if !ok {
				continue
			}

			if instance["instance_number"] == instanceNumber {
				newInstanceList = append(instanceList[:i], instanceList[i+1:]...)

				break
			}
			newInstanceList = append(newInstanceList, inst)
		}
		spec["instance_list"] = newInstanceList

		updatedResource, err := c.dynamicClient.Resource(schema.GroupVersionResource{
			Group:    "oakestra.oakestra.kubernetes",
			Version:  "v1",
			Resource: "oakestrajobs",
		}).Namespace(c.resourceNamespace).Update(context.TODO(), &item, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("OakestraUpdate failed: %w", err)
		}
		log.Printf("Resource updated successfully: %v\n", updatedResource)
	}

	return nil
}
