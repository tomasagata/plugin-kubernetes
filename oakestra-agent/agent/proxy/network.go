package kubenetesproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	k8sclient "oakestra/plugin-kubernetes/oakestra-agent/agent/kubernetesClient"
)

type network struct {
	serviceManagerURL  string
	serviceManagerPort string
}

func (n network) NotifyNetworkDeployment(job k8sclient.OakestraJob) error {
	fmt.Println("Sending network deployment notification to the network component")

	if n.serviceManagerURL == "" || n.serviceManagerPort == "" {
		return fmt.Errorf("cluster Service Manager Service not existent, please restart Agent")
	}

	serviceManagerAddr := "http://" + n.serviceManagerURL + ":" + n.serviceManagerPort + "/api/net/deployment"
	fmt.Println(serviceManagerAddr)
	payload, err := json.Marshal(map[string]string{
		"job_name": job.Spec.JobName,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(serviceManagerAddr, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("calling Service Manager /api/net/deployment not successful: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from Service Manager: %d", resp.StatusCode)
	}

	fmt.Println("Calling Service Manager /api/net/deployment successful.")
	return nil
}
