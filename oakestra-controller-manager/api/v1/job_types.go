/*
Copyright 2024 Jakob Kempter.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OakestraJobSpec defines the desired state of OakestraJob
// serviceSpec defines the desired state of a service in the cluster.
type OakestraJobSpec struct {

	//
	// Specs sent by the Root when requesting the deployment
	//

	// Full Name of Job
	JobName string `json:"job_name"`

	// Contains additional files associated with the service.
	AddedFiles []string `json:"added_files,omitempty"`

	// Unique ID of application
	ApplicationID string `json:"application_ID"`

	// Name of application
	ApplicationName string `json:"application_name"`

	// Namespace of application
	ApplicationNamespace string `json:"application_namespace"`

	// Specifies the incoming bandwidth for the service.
	BandwidthIn int `json:"bandwidth_in,omitempty"`

	// Specifies the outgoing bandwidth for the service.
	BandwidthOut int `json:"bandwidth_out,omitempty"`

	// Commands for the service.
	Cmd []string `json:"cmd,omitempty"`

	// Code which is used by the service
	Code string `json:"code"`

	// Image which is used by the container
	Image string `json:"image"`

	// Specifies the amount of memory needed by the service.
	Memory int `json:"memory,omitempty"`

	// Unique ID of service
	MicroserviceID string `json:"microservice_ID"`

	// Namespace of service
	MicroserviceName string `json:"microservice_name"`

	// Namespace of service
	MicroserviceNamepace string `json:"microservice_namespace"`

	// Instance Number of next instance
	NextInstanceProgressiveNumber int `json:"next_instance_progressive_number"`

	// Specifies the ports used by the service.
	Port string `json:"port"`

	// Represents the current status of the service.
	State string `json:"state,omitempty"`

	// storage specifies the storage used by the service.
	Storage int `json:"storage,omitempty"` // Storage used by the service

	// Represents the number of virtual CPUs used by the service.
	VCPUs int `json:"vcpus,omitempty"`

	// Represents the number of virtual GPUs used by the service.
	VGPUs int `json:"vgpus,omitempty"`

	// Runtime indicates the type of virtualization used by the service.
	Virtualization string `json:"virtualization,omitempty"`

	// Represents the number of virtual TPUs used by the service.
	VTPUs int `json:"vtpus,omitempty"`

	// Represents the current status of the service.
	Status string `json:"status,omitempty"`

	// StatusDetail provides additional details about the status of the service.
	StatusDetail string `json:"status_detail,omitempty"` // Additional details about the status

	// StatusDetail provides additional details about the status of the service.
	InstanceList []Instance `json:"instance_list"` // Additional details about the status

	// Represents the disk size what is needed
	Disk int `json:"disk,omitempty"`

	//
	// Optional Specs
	//

	// // Represents the SLA violation strategy for the service.
	// SlaViolationStrategy string `json:"sla_violation_strategy,omitempty"`

	// // Represents the target node for the service.
	// TargetNode string `json:"target_node,omitempty"`

	// // Contains additional arguments for the service.
	// --group oakestra --version v1 --kind OakestraJob []string `json:"args,omitempty"`

	// // Env contains the environment variables for the service.
	// Env []string `json:"environment,omitempty"`

	// // UnikernelImages contains the list of unikernel images used by the service.
	// UnikernelImages []string `json:"vm_images,omitempty"`

	// // Architectures contains the architectures supported by the service.
	// Architectures []string `json:"arch,omitempty"`

}

// OakestraJobStatus defines the observed state of OakestraJob
type OakestraJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	InstanceList InstanceNummberSet `json:"instanceList"`
}

type Instance struct {
	InstanceNumber  int    `json:"instance_number"`
	ClusterID       string `json:"cluster_ID"`
	ClusterLocation string `json:"cluster_location"`

	Status       string `json:"status,omitempty"`        // reset after schedulement
	StatusDetail string `json:"status_detail,omitempty"` // reset after schedulement

	CPU                   int    `json:"cpu,omitempty"`                     // set after schedulement
	Memory                int    `json:"memory,omitempty"`                  // set after schedulement
	Disk                  int    `json:"disk,omitempty"`                    // set after schedulement
	LastModifiedTimestamp string `json:"last_modified_timestamp,omitempty"` // set after schedulement
	HostPort              string `json:"host_port,omitempty"`               // set after schedulement
	HostIP                string `json:"host_IP,omitempty"`                 // set after schedulement
	WorkerID              string `json:"worker_ID,omitempty"`               // set after schedulement
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OakestraJob is the Schema for the OakestraJobs API
type OakestraJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OakestraJobSpec   `json:"spec,omitempty"`
	Status OakestraJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OakestraJobList contains a list of OakestraJob
type OakestraJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OakestraJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OakestraJob{}, &OakestraJobList{})
}

type InstanceNummberSet map[string]Instance

func (set InstanceNummberSet) Add(info Instance) {
	set[strconv.Itoa(info.InstanceNumber)] = info
}

func (set InstanceNummberSet) Remove(num string) {
	delete(set, num)
}

func (set InstanceNummberSet) Contains(num string) bool {
	_, exists := set[num]
	return exists
}
