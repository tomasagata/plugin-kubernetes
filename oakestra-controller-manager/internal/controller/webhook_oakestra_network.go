package controller

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;delete,versions=v1,name=oakestra.kubernetes.network,admissionReviewVersions=v1,sideEffects=NoneOnDryRun

type WebhookOakestraNetwork struct {
	client.Client
	Decoder     *admission.Decoder
	ClusterInfo ClusterInfo
}

type ClusterInfo struct {
	ClusterID             string
	RootServiceManagerURL string
}

func (oaknet *WebhookOakestraNetwork) undeployInstance(systemJobID string) {
	instanceNumber := "0"
	url := oaknet.ClusterInfo.RootServiceManagerURL + "/api/net/" + systemJobID + "/" + instanceNumber
	sendDeleteRequest(url)
}

func (oaknet *WebhookOakestraNetwork) unregisterService(systemJobID string) {
	url := oaknet.ClusterInfo.RootServiceManagerURL + "/api/net/service/" + systemJobID
	sendDeleteRequest(url)

}

func (oaknet *WebhookOakestraNetwork) deployInstance(systemJobID string) {
	data := SystemJob{
		ClusterID:      oaknet.ClusterInfo.ClusterID,
		InstanceNumber: 0,
		SystemJobID:    systemJobID,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	url := oaknet.ClusterInfo.RootServiceManagerURL + "/api/net/instance/deploy"
	sendPostRequestToRootServiceManager(jsonData, url)
}

func (oaknet *WebhookOakestraNetwork) connectToOakestraNetwork(pod *corev1.Pod) {

	system_job_id_unhashed := pod.Name
	hash := sha256.Sum256([]byte(system_job_id_unhashed))
	system_job_idHashed := hex.EncodeToString(hash[:])[:24]
	systemJobID := system_job_idHashed

	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Annotations["k8s.v1.cni.cncf.io/networks"] = "oakestra-cni"
	pod.Annotations["systemJobID"] = systemJobID

	oaknet.registerService(systemJobID, pod)
	pod.Annotations["oakestra.io/status"] = "registered"

	oaknet.deployInstance(systemJobID)
	pod.Annotations["oakestra.io/status"] = "deployed"
}

func (oaknet *WebhookOakestraNetwork) disconnectFromOakestraNetwork(podName string) {

	systemJobIDUnhashed := podName
	hash := sha256.Sum256([]byte(systemJobIDUnhashed))
	system_job_idHashed := hex.EncodeToString(hash[:])[:24]
	systemJobID := system_job_idHashed

	oaknet.undeployInstance(systemJobID)
	oaknet.unregisterService(systemJobID)
}

func (oaknet *WebhookOakestraNetwork) Handle(ctx context.Context, req admission.Request) admission.Response {
	log := log.FromContext(ctx)

	if req.DryRun != nil && *req.DryRun {
		return admission.Allowed("DryRun requested, admission allowed")
	}

	switch req.Operation {
	case v1.Delete:
		log.Info("DELETE PROCESS \n")
		oaknet.disconnectFromOakestraNetwork(req.Name)
		return admission.Allowed("Deletion allowed")

	case v1.Create:
		log.Info("CREATE PROCESS \n")
		pod := &corev1.Pod{}
		err := oaknet.Decoder.Decode(req, pod)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		hasOakestraNetworkPort, err := oaknet.hasOakestraPortLabel(pod)
		if err != nil {
			admission.Errored(http.StatusInternalServerError, err)
		}
		if hasOakestraNetworkPort {
			oaknet.connectToOakestraNetwork(pod)
		}

		marshaledPod, err := json.Marshal(pod)
		if err != nil {
			return admission.Errored(http.StatusInternalServerError, err)
		}
		return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
	}

	return admission.Allowed("Default")

}

func (a *WebhookOakestraNetwork) hasOakestraPortLabel(pod *corev1.Pod) (bool, error) {
	portValue, portExists := pod.Labels["oakestra.io/port"]
	if !portExists || portValue == "" {
		return false, fmt.Errorf("label oakestra.io/port is not set or empty")
	}
	return true, nil
}

func sendPostRequestToRootServiceManager(jsonData []byte, url string) {

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Response Body:", string(body))
}

func (a *WebhookOakestraNetwork) registerService(systemJobID string, pod *corev1.Pod) {

	descriptor := DeploymentDescriptorNetwork{
		ApplicationID:    systemJobID,
		AppName:          pod.Name,
		AppNamespace:     pod.Namespace,
		ServiceName:      pod.Name,
		ServiceNamespace: pod.Namespace,
		Image:            pod.Spec.Containers[0].Image,

		InstanceList:                  []Instance{},
		Virtualization:                "docker",
		Memory:                        0,
		Storage:                       0,
		NextInstanceProgressiveNumber: 0,
	}

	data := RequestData{
		SystemJobID:          systemJobID,
		DeploymentDescriptor: descriptor,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	url := a.ClusterInfo.RootServiceManagerURL + "/api/net/service/deploy"
	sendPostRequestToRootServiceManager(jsonData, url)
}

func sendDeleteRequest(url string) {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil) // Kein Body gesendet
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Response Body:", string(body))
}

type DeploymentDescriptorNetwork struct {
	ApplicationID    string     `json:"applicationID"`
	AppName          string     `json:"app_name"`
	AppNamespace     string     `json:"app_ns"`
	ServiceName      string     `json:"service_name"`
	ServiceNamespace string     `json:"service_ns"`
	Image            string     `json:"image"`
	InstanceList     []Instance `json:"instance_list"`

	Virtualization                string `json:"virtualization"`
	Memory                        int    `json:"memory"`
	Storage                       int    `json:"storage"`
	NextInstanceProgressiveNumber int    `json:"next_instance_progressive_number"`

	// TODO Future Work.
	// RRIP             string     `json:"RR_ip"`
	// RRIPV6           string     `json:"RR_ip_v6"`
}

type Instance struct{}

type RequestData struct {
	SystemJobID          string                      `json:"system_job_id"`
	DeploymentDescriptor DeploymentDescriptorNetwork `json:"deployment_descriptor"`
}

type SystemJob struct {
	SystemJobID    string `json:"system_job_id"`
	InstanceNumber int    `json:"instance_number"`
	ClusterID      string `json:"cluster_id"`
}
