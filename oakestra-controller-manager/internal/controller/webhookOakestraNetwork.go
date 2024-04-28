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

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update;delete,versions=v1,name=oakestra.kubernetes.network,admissionReviewVersions=v1,sideEffects=None

type WebhookOakestraNetwork struct {
	client.Client
	Decoder    *admission.Decoder
	custerInfo clusterInfo
}

type clusterInfo struct {
	clusterID             string
	rootServiceManagerURL string
}

func (info *clusterInfo) IsSet() bool {
	return info.clusterID != "" && info.rootServiceManagerURL != ""
}

func (a *WebhookOakestraNetwork) Handle(ctx context.Context, req admission.Request) admission.Response {
	log := log.FromContext(ctx)

	// check wether clusterdata is set.

	updateClusterInfo := func() error {
		cm := &corev1.ConfigMap{}
		err := a.Client.Get(ctx, client.ObjectKey{
			Name:      "oakestra-cluster-info",
			Namespace: "oakestra-controller-manager",
		}, cm)
		if err != nil {
			log.Info("ERROR")
			log.Info(err.Error())
			return err
		}

		a.custerInfo.clusterID = cm.Data["CLUSTER_ID"]
		a.custerInfo.rootServiceManagerURL = "http://" + cm.Data["OAKESTRA_ROOT_IP"] + ":" + cm.Data["OAKESTRA_ROOT_NETWORK_PORT"]

		return nil
	}

	hasOakestraNetworkandPortAnnotation := func(pod *corev1.Pod) (bool, error) {

		networkActive := pod.Annotations["oakestraNetwork"] == "active"
		_, portExists := pod.Annotations["oakestraPort"]
		if networkActive && !portExists {
			return false, fmt.Errorf("oakestraPort not set")
		}

		return networkActive, nil

	}

	registerService := func(systemJobID string, pod *corev1.Pod) {

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

		url := a.custerInfo.rootServiceManagerURL + "/api/net/service/deploy"
		log.Info(url)
		sendPostRequestToRootServiceManager(jsonData, url)
	}

	undeployInstance := func(systemJobID string) {
		instanceNumber := "0"
		url := a.custerInfo.rootServiceManagerURL + "/api/net/" + systemJobID + "/" + instanceNumber
		log.Info(url)
		sendDeleteRequest(url)
	}

	unregisterService := func(systemJobID string) {
		url := a.custerInfo.rootServiceManagerURL + "/api/net/service/" + systemJobID
		log.Info(url)
		sendDeleteRequest(url)
	}

	deployInstance := func(systemJobID string) {
		data := SystemJob{
			ClusterID:      a.custerInfo.clusterID,
			InstanceNumber: 0,
			SystemJobID:    systemJobID,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}

		url := a.custerInfo.rootServiceManagerURL + "/api/net/instance/deploy"
		log.Info(url)
		sendPostRequestToRootServiceManager(jsonData, url)
	}

	disconnectFromOakestraNetwork := func(podName string) {

		system_job_id_unhashed := podName
		hash := sha256.Sum256([]byte(system_job_id_unhashed))
		system_job_idHashed := hex.EncodeToString(hash[:])[:24]
		systemJobID := system_job_idHashed

		undeployInstance(systemJobID)
		unregisterService(systemJobID)
	}

	connectToOakestraNetwork := func(pod *corev1.Pod) {

		system_job_id_unhashed := pod.Name
		hash := sha256.Sum256([]byte(system_job_id_unhashed))
		system_job_idHashed := hex.EncodeToString(hash[:])[:24]
		systemJobID := system_job_idHashed

		if pod.Annotations == nil {
			pod.Annotations = map[string]string{}
		}
		pod.Annotations["k8s.v1.cni.cncf.io/networks"] = "oakestra-cni"
		pod.Annotations["systemJobID"] = systemJobID

		registerService(systemJobID, pod)
		pod.Annotations["oakestra-network-status"] = "registered"

		deployInstance(systemJobID)
		pod.Annotations["oakestra-network-status"] = "deployed"
	}

	switch req.Operation {
	case v1.Delete:
		log.Info("DELETE PROCESS \n")

		// TODO - Nur aufrufen, wenn es im Oakestra Netzwerk war!

		disconnectFromOakestraNetwork(req.Name)
		return admission.Allowed("Deletion allowed")
	case v1.Update:
		// disconnectFromOakestraNetwork(req.Name)

		return admission.Allowed("Update allowed")

	case v1.Create:
		log.Info("CREATE PROCESS \n")
		pod := &corev1.Pod{}
		err := a.Decoder.Decode(req, pod)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		if !a.custerInfo.IsSet() {
			log.Info("clusterInfo is not set, updating now...")
			if err := updateClusterInfo(); err != nil {
				log.Error(err, "Failed to update clusterInfo")
				return admission.Errored(http.StatusInternalServerError, err)
			}
		} else {
			log.Info("clusterInfo is already set, no update needed")
		}

		oakestraNetwork, err := hasOakestraNetworkandPortAnnotation(pod)
		if err != nil {
			admission.Errored(http.StatusInternalServerError, err)
		}
		if oakestraNetwork {
			log.Info("HAS OAKESTRA=NETWORK LABEL \n")
			connectToOakestraNetwork(pod)
		}

		marshaledPod, err := json.Marshal(pod)
		if err != nil {
			return admission.Errored(http.StatusInternalServerError, err)
		}
		return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
	}

	return admission.Allowed("Default")

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
