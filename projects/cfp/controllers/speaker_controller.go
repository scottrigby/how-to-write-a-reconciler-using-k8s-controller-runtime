/*
Copyright 2022.

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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	talksv1 "github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime/projects/cfp/api/v1"
)

// SpeakerReconciler reconciles a Speaker object
type SpeakerReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	cfpAPI string
}

//+kubebuilder:rbac:groups=talks.kubecon.na,resources=speakers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=talks.kubecon.na,resources=speakers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=talks.kubecon.na,resources=speakers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Speaker object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *SpeakerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Get the speaker Object

	// Set reconciling condition on the speaker object

	// Set a finalizer on the speaker object if not set

	// Check if this a deletion, if yes api call to delete the Speaker
	// clean the finalizer
	// Clean the condition

	// Check if an ID exist in the Status
		// Check if an update make sense
	  	// Make an Api call and check if anything changed on the speaker spec.
			   // If it is not found 
				 		// Set a condition like CNCFSpeakerErrorConditon
						// error and requeue
				// Update and set ready condition
		// Otherwise unset reconciling and set ready condition

	// Make a call to the API to create speaker
	//...

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpeakerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&talksv1.Speaker{}).
		Complete(r)
}