package kubenetesproxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"oakestra/plugin-kubernetes/oakestra-agent/agent/config"
	k8sclient "oakestra/plugin-kubernetes/oakestra-agent/agent/kubernetesClient"
	"strconv"

	"github.com/gorilla/mux"
)

type Server struct {
	k8sClient k8sclient.KubernetesClient
	config    *config.Config
	network   network
}

func NewServer(k8sClient k8sclient.KubernetesClient, config *config.Config) *Server {
	return &Server{
		k8sClient: k8sClient,
		config:    config,
		network: network{
			serviceManagerURL:  config.ClusterServiceManagerIP,
			serviceManagerPort: strconv.Itoa(config.NetworkComponentPort)},
	}
}

func (s *Server) StartHTTPServer() {
	router := mux.NewRouter()

	router.HandleFunc("/status", s.statusHandler).Methods("GET")
	router.HandleFunc("/api/deploy/{system_job_id}/{instance_number}", s.deployInstanceHandler).Methods("POST")
	router.HandleFunc("/api/delete/{system_job_id}/{instance_number}", s.deleteInstanceHandler).Methods("GET")

	log.Printf("Server is running on port %d...\n", s.config.MyPort)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.config.MyPort), router))
}

func (s *Server) statusHandler(w http.ResponseWriter, _ *http.Request) {
	log.Println("Incoming Request /status")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok\n")
}

func (s *Server) deployInstanceHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Incoming Request /api/deploy")

	vars := mux.Vars(r)
	systemJobID := vars["system_job_id"]

	instanceNumber, err := strconv.Atoi(vars["instance_number"])
	if err != nil {
		log.Printf("Could not convert String to Integer: %s", err)
		return
	}

	log.Printf("system_job_id: %s, instance_number: %d\n", systemJobID, instanceNumber)

	var requestStruct k8sclient.OakestraJSONRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&requestStruct)
	if err != nil {
		log.Printf("Could not parse body to Struct: %s", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	requestStruct.InstanceList = filterInstances(requestStruct.InstanceList, s.config.ClusterID)

	oakestraJob, err := s.k8sClient.CreateOakestraJob(requestStruct, systemJobID, instanceNumber)
	if err != nil {
		log.Printf("Failed to create OakestraJob: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create OakestraJob: %v", err), http.StatusInternalServerError)
		return
	}

	err = s.network.NotifyNetworkDeployment(oakestraJob)
	if err != nil {
		log.Printf("Error notifying network: %s", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok\n")
}

func filterInstances(input []k8sclient.InstanceJSONRequest, clusterID string) []k8sclient.InstanceJSONRequest {
	result := []k8sclient.InstanceJSONRequest{}

	for _, instance := range input {
		if instance.ClusterID == clusterID {
			result = append(result, instance)
		}
	}
	return result
}

func (s *Server) deleteInstanceHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Incoming Request /api/delete/ - to delete task...")

	vars := mux.Vars(r)
	systemJobID := vars["system_job_id"]

	instanceNumber, err := strconv.ParseInt(vars["instance_number"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting instance number: %v", err), http.StatusBadRequest)
		return
	}

	err = s.k8sClient.DeleteInstance(systemJobID, instanceNumber)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete instance: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok\n")
}
