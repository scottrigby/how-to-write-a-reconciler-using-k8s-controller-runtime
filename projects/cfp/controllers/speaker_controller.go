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
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/json"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/fluxcd/pkg/apis/meta"
	"github.com/fluxcd/pkg/runtime/conditions"
	"github.com/fluxcd/pkg/runtime/patch"

	talksv1 "github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime/projects/cfp/api/v1"
	"github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime/projects/cfp/internal/cfp"
)

const speakerPath = "/speakers"

var speakerOwnedConditions = []string{
	meta.ReadyCondition,
	meta.ReconcilingCondition,
	meta.StalledCondition,
	talksv1.CreateFailedCondition,
	talksv1.UpdateFailedCondition,
}

// SpeakerReconciler reconciles a Speaker object
type SpeakerReconciler struct {
	client.Client
	HTTPClient     *http.Client
	ControllerName string
	CfpAPI         string
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
	log := log.FromContext(ctx)
	// 1. Fetch the Speaker instance
	// Automatically requeue if an error is returned
	// otherwise requeue based on the result.requeue and result.requeueAfter
	obj := &talksv1.Speaker{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("reconciling speaker", "speaker", obj.Name)

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
		// See https://alenkacz.medium.com/kubernetes-operator-best-practices-implementing-observedgeneration-250728868792
		if conditions.IsStalled(obj) || conditions.IsReady(obj) {
			patchOpts = append(patchOpts, patch.WithStatusObservedGeneration{})
		}

		// Finally, patch the resource
		if err := patchHelper.Patch(ctx, obj, patchOpts...); err != nil {
			if !obj.GetDeletionTimestamp().IsZero() {
				err = kerrors.FilterOut(err, func(e error) bool { return apierrors.IsNotFound(e) })
			}

			retErr = kerrors.NewAggregate([]error{retErr, err})
		}
	}()

	// 4. Set a finalizer on the obj object if not set
	// Finalizers are keys on resources that signal pre-delete operations. They control the garbage collection on resources,
	// and are designed to alert controllers what cleanup operations to perform prior to removing a resource.
	// https://kubernetes.io/blog/2021/05/14/using-finalizers-to-control-deletion/
	if !controllerutil.ContainsFinalizer(obj, talksv1.Finalizer) {
		controllerutil.AddFinalizer(obj, talksv1.Finalizer)
		return ctrl.Result{Requeue: true}, nil
	}

	//create a new speaker client
	speakerClient := cfp.NewSpeakerClient(r.CfpAPI, r.HTTPClient)

	// 5. Check if this a deletion, if yes api call to delete the Speaker
	if !obj.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, obj, speakerClient)
	}

	// 6. Perform reconciliation logic
	result, retErr = r.reconcile(ctx, obj, speakerClient)
	return
}

// reconcile will perform the reconciliation logic for the Speaker object
// While reconciling if an error is encountered, it sets the failure details  in the appropriate
// status condition and returns the error.
// Step 6
func (r *SpeakerReconciler) reconcile(ctx context.Context, obj *talksv1.Speaker, client *cfp.SpeakerClient) (result ctrl.Result, retErr error) {
	// Step 6.3
	// defer func attempt to set the Ready condition and unset all needed conditions based on the reconciliation
	defer func() {
		if !result.Requeue && retErr == nil {
			conditions.Delete(obj, meta.ReconcilingCondition)
			conditions.Delete(obj, talksv1.CreateFailedCondition)
			conditions.Delete(obj, talksv1.UpdateFailedCondition)
			conditions.MarkTrue(obj, meta.ReadyCondition, meta.SucceededReason, "reconciled '%s' successfully", obj.Name)
		}

		if retErr != nil {
			var apiErr *cfp.Error
			if ok := errors.As(retErr, &apiErr); !ok {
				apiErr = &cfp.Error{
					Reason: cfp.ErrorReason{
						Reason:  "unknown",
						Summary: "unkonwn error",
					},
					Err: retErr,
				}
			}
			switch apiErr.Reason {
			case cfp.ErrCreateSpeaker:
				conditions.MarkTrue(obj, talksv1.CreateFailedCondition, apiErr.Reason.Reason, apiErr.Error())
				conditions.MarkFalse(obj, meta.ReadyCondition, apiErr.Reason.Reason, apiErr.Error())
			case cfp.ErrUpdateSpeaker:
				conditions.MarkTrue(obj, talksv1.UpdateFailedCondition, apiErr.Reason.Reason, apiErr.Error())
				conditions.MarkFalse(obj, meta.ReadyCondition, apiErr.Reason.Reason, apiErr.Error())
			case cfp.ErrCreateRequest, cfp.ErrMakeRequest, cfp.ErrFetchSpeaker:
				conditions.MarkFalse(obj, meta.ReadyCondition, apiErr.Reason.Reason, apiErr.Error())
			default:
				conditions.MarkFalse(obj, meta.ReadyCondition, meta.FailedReason, apiErr.Error())
			}
		}
	}()

	// Set the initial status of the object to be reconciling
	if obj.Generation != obj.Status.ObservedGeneration {
		conditions.MarkReconciling(obj, meta.ProgressingReason, fmt.Sprintf("Reconciling a new generation of the object %d", obj.Generation))
	}

	// Step 6.1
	// Check if an ID exist in the Status
	// Check if an update make sense
	// Make an Api call and check if anything changed on the obj spec.
	// If it is not found
	// Set a condition UpdateFailedCondition
	// error and requeue
	if obj.Status.ID != "" {
		err := r.handleSpeakerUpdate(ctx, obj, client)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Step 6.2
	// Create the Speaker
	err := r.createSpeaker(ctx, obj, client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Set the ID in the status
	obj.Status.ID = fmt.Sprintf("%s-%s", obj.Namespace, obj.Name)
	return ctrl.Result{}, nil
}

func (r *SpeakerReconciler) handleSpeakerUpdate(ctx context.Context, obj *talksv1.Speaker, client *cfp.SpeakerClient) error {
	// Make a call to the API to update obj
	body, err := createpayload(obj)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Get the Speaker
	payload, err := client.Get(ctx, speakerPath, obj.Status.ID)
	if err != nil {
		return err
	}

	// Check if the payload is the same
	// if not, it means we have changed the spec, so we need to update the speaker
	if !bytes.Equal(payload, body) {
		err = client.Update(ctx, speakerPath, obj.Status.ID, body)
		if err != nil {
			return err
		}
	}

	return nil
}

// createSpeaker will create a Speaker in the CFP API
func (r *SpeakerReconciler) createSpeaker(ctx context.Context, obj *talksv1.Speaker, client *cfp.SpeakerClient) error {
	// Make a call to the API to update obj
	body, err := createpayload(obj)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	err = client.Create(ctx, speakerPath, body)
	if err != nil {
		return err
	}

	return nil
}

// reconcileDelete will delete the obj from the CFP API
// Step 5
func (r *SpeakerReconciler) reconcileDelete(ctx context.Context, obj *talksv1.Speaker, client *cfp.SpeakerClient) (ctrl.Result, error) {
	// api call to delete the Speaker
	err := client.Delete(ctx, speakerPath, obj.Status.ID)
	if err != nil {
		// return the error so we can requeue
		return ctrl.Result{}, err
	}
	// clean the finalizer
	controllerutil.RemoveFinalizer(obj, talksv1.Finalizer)
	fmt.Println("FINALIZER REMOVED", obj.Generation)

	// Stop the reconciliation
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpeakerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&talksv1.Speaker{}).
		Complete(r)
}

func createpayload(obj *talksv1.Speaker) ([]byte, error) {
	body := struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Bio   string `json:"bio"`
		Email string `json:"email"`
	}{
		ID:    fmt.Sprintf("%s-%s", obj.Namespace, obj.Name),
		Name:  obj.Spec.Name,
		Bio:   obj.Spec.Bio,
		Email: obj.Spec.Email,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshalling body: %w", err)
	}

	return payload, nil
}
