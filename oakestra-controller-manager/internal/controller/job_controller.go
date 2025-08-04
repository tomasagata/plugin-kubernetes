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

package controller

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	oakestrav1 "oakestra/OakestraJob-operator/api/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// OakestraJobReconciler reconciles a OakestraJob object
type OakestraJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type OakestraEmptyInstanceListError struct {
	CRDName string
}

func (e *OakestraEmptyInstanceListError) Error() string {
	return fmt.Sprintf("%s: InstanceList is empty", e.CRDName)
}

// +kubebuilder:rbac:groups=oakestra.oakestra.kubernetes,resources=oakestrajobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=oakestra.oakestra.kubernetes,resources=oakestrajobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=oakestra.oakestra.kubernetes,resources=oakestrajobs/finalizers,verbs=update

//
// Pod Access Control
//
//+kubebuilder:rbac:groups=oakestra.oakestra.kubernetes,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=oakestra.oakestra.kubernetes,resources=pods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=oakestra.oakestra.kubernetes,resources=pods/finalizers,verbs=update

// Pod Access Control
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups="",resources=pods/finalizers,verbs=update

//
// Deployment Access Control
//
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update

//
// ConfigMap Access Control
//
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

func (r *OakestraJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	InstanceList := make(oakestrav1.InstanceNummberSet)

	var oakestraJob oakestrav1.OakestraJob
	if err := r.Get(ctx, req.NamespacedName, &oakestraJob); err != nil {
		log.Error(err, "unable to fetch Oakestra OakestraJob")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	getEnvVariables := func() []corev1.EnvVar {
		var envVars []corev1.EnvVar

		for _, variable := range oakestraJob.Spec.Environment {
			parts := strings.SplitN(variable, "=", 2)
			if len(parts) == 2 {

				envVars = append(envVars, corev1.EnvVar{Name: parts[0], Value: parts[1]})
			} else {
				log.Info("ENV Variable could not be decoded: ", variable)
			}
		}

		return envVars
	}

	updateInstanceInformation := func(instance *oakestrav1.Instance) {
		selector := labels.SelectorFromSet(labels.Set(map[string]string{"MicroserviceID": oakestraJob.Spec.MicroserviceID}))
		pods := &corev1.PodList{}
		if err := r.Client.List(context.TODO(), pods, &client.ListOptions{
			Namespace:     req.Namespace,
			LabelSelector: selector,
		}); err != nil {
			log.Error(err, "Failed to list pods")
			return
		}

		if len(pods.Items) == 0 {
			log.Info("Podlist is empty.")
		} else {

			podPhase := pods.Items[0].Status.Phase
			switch podPhase {
			case corev1.PodPending:
				instance.Status = "NO_WORKER_CAPACITY"
			case corev1.PodRunning:
				instance.Status = "RUNNING"
			case corev1.PodFailed:
				instance.Status = "FAILED"
			}

			instance.CPU = oakestraJob.Spec.VCPUs
			instance.Memory = oakestraJob.Spec.Memory
			instance.StatusDetail = oakestraJob.Spec.StatusDetail
			instance.Disk = oakestraJob.Spec.Disk
			instance.HostIP = pods.Items[0].Status.HostIP

			workerID := pods.Items[0].Spec.NodeName + "-" + pods.Items[0].Status.HostIP
			hash := sha256.Sum256([]byte(workerID))
			clientIDHashed := hex.EncodeToString(hash[:])[:24]
			instance.WorkerID = clientIDHashed

			instance.HostPort = "50011"
		}
	}

	setDeploymentPodTemplate := func(deployment *appsv1.Deployment, oakestraJob *oakestrav1.OakestraJob, instanceNumber string) error {

		deployment.ObjectMeta.Labels = map[string]string{
			"applicationName":  oakestraJob.Spec.ApplicationName,
			"applicationID":    oakestraJob.Spec.ApplicationID,
			"microserviceName": oakestraJob.Spec.MicroserviceName,
			"microserviceID":   oakestraJob.Spec.MicroserviceID,
			"instanceNumber":   instanceNumber,
		}

		// Hier verwenden wir die Labels des Deployments anstelle der Labels des Pods
		deployment.Spec.Template.ObjectMeta.Labels = map[string]string{
			"applicationName":  oakestraJob.Spec.ApplicationName,
			"applicationID":    oakestraJob.Spec.ApplicationID,
			"microserviceName": oakestraJob.Spec.MicroserviceName,
			"microserviceID":   oakestraJob.Spec.MicroserviceID,
			"instanceNumber":   instanceNumber,
		}

		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"applicationName":  oakestraJob.Spec.ApplicationName,
				"applicationID":    oakestraJob.Spec.ApplicationID,
				"microserviceName": oakestraJob.Spec.MicroserviceName,
				"microserviceID":   oakestraJob.Spec.MicroserviceID,
				"instanceNumber":   instanceNumber,
			},
		}

		deployment.Spec.Template.Annotations = map[string]string{
			"k8s.v1.cni.cncf.io/networks": "oakestra-cni",
			"oakestraPort":                oakestraJob.Spec.Port,
		}

		// replicaCount := int32(10)
		// deployment.Spec.Replicas = &replicaCount
		deployment.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyAlways

		deployment.Spec.Template.Spec.Containers = []corev1.Container{
			{
				Name:    oakestraJob.Spec.MicroserviceName,
				Image:   oakestraJob.Spec.Image,
				Command: oakestraJob.Spec.Cmd,
				Env:     getEnvVariables(),
			},
		}

		return nil
	}

	OakestraJobName := fmt.Sprintf(
		"%s.%s.%s.%s",
		oakestraJob.Spec.ApplicationName,
		oakestraJob.Spec.ApplicationNamespace,
		oakestraJob.Spec.MicroserviceName,
		oakestraJob.Spec.MicroserviceNamepace,
	)

	instances := oakestraJob.Spec.InstanceList
	if len(instances) == 0 {
		log.Info("InstanceList is empty. No Oakestra Pods are deployed")
	}

	updatedInstanceList := []oakestrav1.Instance{}

	//
	// Create missing instances.
	//

	for _, instance := range instances {
		InstanceList.Add(oakestrav1.Instance{InstanceNumber: instance.InstanceNumber})

		deploy := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s.%d", OakestraJobName, instance.InstanceNumber),
				Namespace: req.Namespace},
		}

		_, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, deploy, func() error {
			err := setDeploymentPodTemplate(deploy, &oakestraJob, strconv.Itoa(instance.InstanceNumber))
			if err != nil {
				return err
			}
			if err := controllerutil.SetControllerReference(&oakestraJob, deploy, r.Scheme); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			if strings.Contains(err.Error(), "the object has been modified; please apply your changes to the latest version and try again") {
				log.Info("Deployment has already been modified")
			} else {
				log.Error(err, "Deployment reconcile failed")
			}
		}

		updateInstanceInformation(&instance)
		updatedInstanceList = append(updatedInstanceList, instance)
	}
	oakestraJob.Spec.InstanceList = updatedInstanceList

	//
	// Delete process
	//

	var deployments appsv1.DeploymentList
	if err := r.List(ctx, &deployments, client.InNamespace(oakestraJob.Namespace), client.MatchingFields{OakestraJobOwnerKey: oakestraJob.Name}); err != nil {
		log.Error(err, "unable to get all deployments")
	}

	for _, deployment := range deployments.Items {
		instanceNumberStr, found := deployment.Labels["instanceNumber"]
		if !found {
			log.Error(fmt.Errorf("deployment %s/%s is missing 'instanceNumber' label", deployment.Namespace, deployment.Name), "")
			continue
		}

		if !InstanceList.Contains(instanceNumberStr) {
			log.Info(fmt.Sprintf("Deleting Deployment %s/%s as its instance number %s is not in the InstanceList", deployment.Namespace, deployment.Name, instanceNumberStr))
			if err := r.Delete(ctx, &deployment); err != nil {
				if k8serrors.IsNotFound(err) {
					log.Info("Deployment deleted successfully")
					break
				}
				log.Error(err, "unable to delete Deployment")
				return ctrl.Result{}, err
			}
		}
	}

	// Update Status of OakestraJob
	oakestraJob.Status.InstanceList = InstanceList

	// Update new Changes - for status upgrades of instances
	if err := r.Update(ctx, &oakestraJob); err != nil {
		if strings.Contains(err.Error(), "the object has been modified; please apply your changes to the latest version and try again") {
			log.Info("OakestraJob has already been modified")
			log.Info(err.Error())
		} else {
			log.Error(err, "unable to update Oakestra OakestraJob Spec")
			return ctrl.Result{}, err
		}
	}

	if len(InstanceList) == 0 {
		if err := r.Delete(ctx, &oakestraJob); err != nil {
			log.Error(err, "unable to delete Oakestra OakestraJob")
			return ctrl.Result{}, err
		}
		log.Info("OakestraJob deleted")

	}

	return ctrl.Result{}, nil
}

var (
	OakestraJobOwnerKey = ".metadata.controller"
	apiGVStr            = oakestrav1.GroupVersion.String()
)

// SetupWithManager sets up the controller with the Manager.
func (r *OakestraJobReconciler) SetupWithManager(mgr ctrl.Manager) error {

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.Deployment{}, OakestraJobOwnerKey, func(rawObj client.Object) []string {
		// grab the OakestraJob object, extract the owner...
		deployment := rawObj.(*appsv1.Deployment)
		owner := metav1.GetControllerOf(deployment)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != apiGVStr || owner.Kind != "OakestraJob" {
			return nil
		}

		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&oakestrav1.OakestraJob{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
