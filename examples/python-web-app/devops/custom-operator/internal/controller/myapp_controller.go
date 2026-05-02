/*
Copyright 2026.

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

	// Standard Kubernetes APIs
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	// Controller Runtime
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	// Your Custom API
	webappv1 "github.com/laxmangodi/Laxman_Docker_end_to_end/examples/python-web-app/custom-operator/api/v1"
)

// MyAppReconciler reconciles a MyApp object
type MyAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.example.com,resources=myapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.example.com,resources=myapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.example.com,resources=myapps/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch

func (r *MyAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// 1. Fetch the MyApp instance
	myapp := &webappv1.MyApp{} // Use webappv1 alias
	err := r.Get(ctx, req.NamespacedName, myapp)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// 2. Check if the Deployment already exists
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: myapp.Name, Namespace: myapp.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		// 3. Define and Create a new Deployment
		dep := r.deploymentForMyApp(myapp)
		l.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			l.Error(err, "Failed to create new Deployment")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// 4. Ensure the replica count (size) matches the Spec
	size := myapp.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			l.Error(err, "Failed to update Deployment replicas")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// deploymentForMyApp returns a myapp Deployment object
func (r *MyAppReconciler) deploymentForMyApp(m *webappv1.MyApp) *appsv1.Deployment {
	ls := map[string]string{"app": "myapp", "myapp_cr": m.Name}
	replicas := m.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "laxmangodi/python-web-app:v1",
						Name:  "python-web-app",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "myapp",
						}},
					}},
				},
			},
		},
	}
	// Set MyApp instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.Scheme)
	return dep
}

func (r *MyAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.MyApp{}).     // Use webappv1 alias
		Owns(&appsv1.Deployment{}). // Tell operator to watch deployments it owns
		Complete(r)

}
