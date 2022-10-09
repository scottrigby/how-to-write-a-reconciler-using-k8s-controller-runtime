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
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/fluxcd/pkg/apis/meta"
	"github.com/fluxcd/pkg/runtime/conditions"
	"github.com/fluxcd/pkg/runtime/patch"

	talksv1 "github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime/projects/cfp/api/v1"
)

var speakerOwnedConditions = []string{
	meta.ReadyCondition,
	meta.ReconcilingCondition,
	meta.StalledCondition,
	talksv1.FetchFailedCondition,
}

// SpeakerReconciler reconciles a Speaker object
type SpeakerReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	ControllerName string
	cfpAPI         string
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
func (r *SpeakerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, retErr error) {
	_ = log.FromContext(ctx)

	// 1. Fetch the Speaker instance
	obj := &talksv1.Speaker{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2.Initialize the patch helper with the current version of the object.
	// Helper is a utility for ensuring the proper patching of objects.
	patchHelper, err := patch.NewHelper(obj, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// 3. Always attempt to Patch the Speaker object and status after each reconciliation.
	defer func() {
		// Patch the object, ignoring conflicts on the conditions owned by this controller
		patchOpts := []patch.Option{
			patch.WithOwnedConditions{
				Conditions: speakerOwnedConditions,
			},
		}

		patchOpts = append(patchOpts, patch.WithFieldOwner(r.ControllerName))

		// Set status observed generation field if the object is stalled, or ready.
		if conditions.IsStalled(obj) || conditions.IsReady(obj) {
			patchOpts = append(patchOpts, patch.WithStatusObservedGeneration{})
		}

		// Finally, patch the resource
		if err := patchHelper.Patch(ctx, obj, patchOpts...); err != nil {
			retErr = kerrors.NewAggregate([]error{retErr, err})
		}
	}()

	// 4. Set a finalizer on the speaker object if not set
	// Finalizers are keys on resources that signal pre-delete operations. They control the garbage collection on resources,
	// and are designed to alert controllers what cleanup operations to perform prior to removing a resource.
	// https://kubernetes.io/blog/2021/05/14/using-finalizers-to-control-deletion/
	if !controllerutil.ContainsFinalizer(obj, talksv1.Finalizer) {
		controllerutil.AddFinalizer(obj, talksv1.Finalizer)
		return ctrl.Result{Requeue: true}, nil
	}

	// 5. Check if this a deletion, if yes api call to delete the Speaker
	if !obj.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, obj)
	}

	// 6. Perform reconciliation logic
	result, retErr = r.reconcile(ctx, obj)
	return
}

// reconcile will perform the reconciliation logic for the Speaker object
// While reconciling if an error is encountered, it sets the failure details  in the appropriate
// status condition and returns the error.
// Step 6
func (r *SpeakerReconciler) reconcile(ctx context.Context, obj *talksv1.Speaker) (ctrl.Result, error) {
	// Set the initial status of the object to be reconciling

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

// reconcileDelete will delete the speaker from the CFP API
// Step 5
func (r *SpeakerReconciler) reconcileDelete(ctx context.Context, obj *talksv1.Speaker) (ctrl.Result, error) {
	// api call to delete the Speaker
	// clean the finalizer
	// Clean the condition
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpeakerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&talksv1.Speaker{}).
		Complete(r)
}
