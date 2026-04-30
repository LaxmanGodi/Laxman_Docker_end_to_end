package main

import (
	"context"
	// "fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap" // Industry standard logger
)

// --- 1. DATA STRUCTURES ---

type DatabaseSpec struct {
	Engine    string `json:"engine"`
	StorageGb int    `json:"storageGb"`
}

type Database struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DatabaseSpec `json:"spec"`
}

// Required for Go's internal machinery
func (in *Database) DeepCopyObject() runtime.Object {
	out := *in
	return &out
}

type DatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Database `json:"items"`
}

func (in *DatabaseList) DeepCopyObject() runtime.Object {
	out := *in
	return &out
}

// --- 2. RECONCILER ---

type DatabaseReconciler struct {
	client.Client
}

func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.Log.WithValues("database", req.NamespacedName)

	var db Database
	if err := r.Get(ctx, req.NamespacedName, &db); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Successfully reconciled", "engine", db.Spec.Engine)
	return ctrl.Result{}, nil
}

// --- 3. MAIN ---

func main() {
	// FIX 1: Set the logger so you can see why it crashes
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	setupLog := ctrl.Log.WithName("setup")

	// FIX 2: Create a proper Scheme with Group and Version
	scheme := runtime.NewScheme()
	gv := schema.GroupVersion{Group: "stable.example.com", Version: "v1"}

	// Register both the object and the list type
	scheme.AddKnownTypes(gv, &Database{}, &DatabaseList{})
	metav1.AddToGroupVersion(scheme, gv)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// FIX 3: Setup the controller
	err = ctrl.NewControllerManagedBy(mgr).
		For(&Database{}).
		Complete(&DatabaseReconciler{Client: mgr.GetClient()})

	if err != nil {
		setupLog.Error(err, "unable to create controller")
		os.Exit(1)
	}

	setupLog.Info("🚀 Industry Standard Go Controller is starting...")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
