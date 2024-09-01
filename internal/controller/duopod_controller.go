/*
Copyright 2024.

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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "github.com/example/duopod-operator/api/v1"
)

// DuoPodReconciler reconciles a DuoPod object
type DuoPodReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.example.com,resources=duopods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.example.com,resources=duopods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.example.com,resources=duopods/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete

func (r *DuoPodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the DuoPod instance
	duopod := &appsv1.DuoPod{}
	err := r.Get(ctx, req.NamespacedName, duopod)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("DuoPod resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get DuoPod")
		return ctrl.Result{}, err
	}

	// Check if the pods already exist, if not create them
	podNames := []string{
		fmt.Sprintf("%s-pod-1", duopod.Name),
		fmt.Sprintf("%s-pod-2", duopod.Name),
	}

	for _, podName := range podNames {
		found := &corev1.Pod{}
		err = r.Get(ctx, types.NamespacedName{Name: podName, Namespace: duopod.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			// Define a new pod
			pod := r.podForDuoPod(duopod, podName)
			log.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
			err = r.Create(ctx, pod)
			if err != nil {
				log.Error(err, "Failed to create new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
				return ctrl.Result{}, err
			}
		} else if err != nil {
			log.Error(err, "Failed to get Pod")
			return ctrl.Result{}, err
		}
	}

	// Update the DuoPod status
	duopod.Status.PodNames = podNames
	err = r.Status().Update(ctx, duopod)
	if err != nil {
		log.Error(err, "Failed to update DuoPod status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// podForDuoPod returns a pod for the DuoPod
func (r *DuoPodReconciler) podForDuoPod(m *appsv1.DuoPod, podName string) *corev1.Pod {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: m.Namespace,
			Labels:    map[string]string{"app": "duopod", "duopod": m.Name},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{m.Spec.Container},
		},
	}
	err := ctrl.SetControllerReference(m, pod, r.Scheme)
	if err != nil {
		return nil
	}
	return pod
}

// SetupWithManager sets up the controller with the Manager.
func (r *DuoPodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.DuoPod{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
