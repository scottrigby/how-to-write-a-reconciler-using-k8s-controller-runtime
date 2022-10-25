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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/fluxcd/pkg/apis/meta"
	"github.com/fluxcd/pkg/runtime/conditions"
	"github.com/fluxcd/pkg/runtime/patch"
	talksv1 "github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime/cfp/api/v1"
	"github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime/cfp/internal/cfp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var proposalOwnedConditions = []string{
	meta.ReadyCondition,
	meta.ReconcilingCondition,
	meta.StalledCondition,
	talksv1.CreateFailedCondition,
	talksv1.UpdateFailedCondition,
	talksv1.FetchFailedCondition,
}

// ProposalReconciler reconciles a Proposal object
type ProposalReconciler struct {
	client.Client
	HTTPClient     *http.Client
	ControllerName string
	CfpAPI         string
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProposalReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetCache().IndexField(context.TODO(), &talksv1.Proposal{}, talksv1.SpeakerIndexKey,
		r.indexProposalBySpeakerName); err != nil {
		return fmt.Errorf("failed setting index fields: %w", err)
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&talksv1.Proposal{}).
		Watches(
			&source.Kind{Type: &talksv1.Speaker{}},
			handler.EnqueueRequestsFromMapFunc(r.requestsForSpeakerChange),
			builder.WithPredicates(SpeakerChangePredicate{}),
		).
		Complete(r)
}

func (r *ProposalReconciler) indexProposalBySpeakerName(o client.Object) []string {
	p, ok := o.(*talksv1.Proposal)
	if !ok {
		panic(fmt.Sprintf("Expected a Proposal, got %T", o))
	}
	return []string{p.Spec.SpeakerRef.Name}
}

func (r *ProposalReconciler) requestsForSpeakerChange(o client.Object) []reconcile.Request {
	speaker, ok := o.(*talksv1.Speaker)
	if !ok {
		panic(fmt.Sprintf("Expected a Speaker, got %T", o))
	}
	// If we do not have an ID, we can't look up proposals
	if speaker.Status.ID == "" {
		return nil
	}

	ctx := context.Background()
	var list talksv1.ProposalList
	if err := r.List(ctx, &list, client.MatchingFields{talksv1.SpeakerIndexKey: speaker.Name}); err != nil {
		return nil
	}

	var reqs []reconcile.Request
	for _, i := range list.Items {
		reqs = append(reqs, reconcile.Request{NamespacedName: client.ObjectKeyFromObject(&i)})
	}
	return reqs
}

//+kubebuilder:rbac:groups=talks.kubecon.na,resources=proposals,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=talks.kubecon.na,resources=proposals/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=talks.kubecon.na,resources=proposals/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Proposal object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ProposalReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, retErr error) {
	log := log.FromContext(ctx)

	// Fetch the proposal
	// Automatically requeue if an error is returned
	// otherwise requeue based on the result.requeue and result.requeueAfter
	obj := &talksv1.Proposal{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("reconciling proposal", "speaker", obj.Name)

	// Initialize the patch helper with the current version of the object.
	// Helper is a utility for ensuring the proper patching of objects.
	patchHelper, err := patch.NewHelper(obj, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always attempt to Patch the Proposal object and status after each reconciliation.
	defer func() {
		// Patch the object, ignoring conflicts on the conditions owned by this controller
		patchOpts := []patch.Option{
			patch.WithOwnedConditions{
				Conditions: proposalOwnedConditions,
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

	// Set a finalizer on the obj object if not set
	// Finalizers are keys on resources that signal pre-delete operations. They control the garbage collection on resources,
	// and are designed to alert controllers what cleanup operations to perform prior to removing a resource.
	// https://kubernetes.io/blog/2021/05/14/using-finalizers-to-control-deletion/
	if !controllerutil.ContainsFinalizer(obj, talksv1.Finalizer) {
		controllerutil.AddFinalizer(obj, talksv1.Finalizer)
		return ctrl.Result{Requeue: true}, nil
	}

	//create a new propsal client
	cfpClient, err := cfp.NewClient(r.CfpAPI, r.HTTPClient)
	if err != nil {
		conditions.MarkStalled(obj, talksv1.CreateFailedCondition, "Failed to create CFP client")
		return ctrl.Result{}, err
	}

	// Check if this a deletion, if yes api call to delete the Speaker
	if !obj.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, obj, cfpClient)
	}

	// Perform reconciliation logic
	result, retErr = r.reconcile(ctx, obj, cfpClient)

	return
}

func (r *ProposalReconciler) reconcile(ctx context.Context, obj *talksv1.Proposal, client *cfp.Client) (result ctrl.Result, retErr error) {
	// defer func attempt to set the Ready condition and unset all needed conditions based on the reconciliation
	defer func() {
		if !result.Requeue && retErr == nil {
			conditions.Delete(obj, meta.ReconcilingCondition)
			conditions.Delete(obj, talksv1.CreateFailedCondition)
			conditions.Delete(obj, talksv1.FetchFailedCondition)
			conditions.Delete(obj, talksv1.UpdateFailedCondition)
			conditions.MarkTrue(obj, meta.ReadyCondition, meta.SucceededReason, "reconciled '%s' successfully", obj.Name)
		}

		if retErr != nil {
			var apiErr *cfp.Error
			if ok := errors.As(retErr, &apiErr); ok {
				switch apiErr.Reason {
				case cfp.ErrCreateProposal:
					conditions.MarkTrue(obj, talksv1.CreateFailedCondition, apiErr.Reason.Reason, apiErr.Error())
					conditions.MarkFalse(obj, meta.ReadyCondition, apiErr.Reason.Reason, apiErr.Error())
				case cfp.ErrUpdateProposal:
					conditions.MarkTrue(obj, talksv1.UpdateFailedCondition, apiErr.Reason.Reason, apiErr.Error())
					conditions.MarkFalse(obj, meta.ReadyCondition, apiErr.Reason.Reason, apiErr.Error())
				case cfp.ErrCreateRequest, cfp.ErrMakeRequest, cfp.ErrFetchProposal:
					conditions.MarkFalse(obj, meta.ReadyCondition, apiErr.Reason.Reason, apiErr.Error())
				default:
					conditions.MarkFalse(obj, meta.ReadyCondition, meta.FailedReason, apiErr.Error())
				}
			}
		}
	}()

	// Set the initial status of the object to be reconciling
	if obj.Generation != obj.Status.ObservedGeneration {
		conditions.MarkReconciling(obj, meta.ProgressingReason, fmt.Sprintf("Reconciling a new generation of the object %d", obj.Generation))
	}

	// Get the speaker reference
	speaker := &talksv1.Speaker{}
	namespacedName := types.NamespacedName{Namespace: obj.Namespace, Name: obj.Spec.SpeakerRef.Name}
	if obj.Spec.SpeakerRef.Namespace != "" {
		namespacedName.Namespace = obj.Spec.SpeakerRef.Namespace
	}

	if err := r.Get(ctx, namespacedName, speaker); err != nil {
		err = fmt.Errorf("unable to get speaker %s: %w", namespacedName.String(), err)
		conditions.MarkTrue(obj, talksv1.FetchFailedCondition, talksv1.FetchFailedCondition, err.Error())
		conditions.MarkFalse(obj, meta.ReadyCondition, meta.FailedReason, err.Error())
		return ctrl.Result{}, err
	}

	if speaker.Status.ID == "" {
		err := fmt.Errorf("unable to get speaker %s", namespacedName.String())
		conditions.MarkTrue(obj, talksv1.FetchFailedCondition, talksv1.FetchFailedCondition, err.Error())
		conditions.MarkFalse(obj, meta.ReadyCondition, meta.FailedReason, err.Error())
		return ctrl.Result{}, err
	}

	speakerID := speaker.Status.ID

	// If we have a Submission on the ProposalStatus sub resource
	// Check if an update is needed
	// If the proposal is marked final, and the submission status is not final, create an entry in cfp.
	var (
		response *ProposalObject
		err      error
	)
	if obj.Status.Submission != "" {
		response, err = r.updateSubmission(ctx, obj, speakerID, client)
		if err != nil {
			return ctrl.Result{}, err
		}

		if response != nil {
			obj.Status.Submission = response.Submission.Status
			obj.Status.LastUpdate = metav1.Time{Time: response.Submission.LastUpdate}
		}
		return ctrl.Result{}, nil
	}

	// Create the Proposal
	response, err = r.createProposal(ctx, obj, speakerID, client)
	if err != nil {
		return ctrl.Result{}, err
	}

	obj.Status.Submission = response.Submission.Status
	obj.Status.LastUpdate = metav1.Time{Time: response.Submission.LastUpdate}

	return ctrl.Result{}, nil
}

func (r *ProposalReconciler) createProposal(ctx context.Context, obj *talksv1.Proposal, speakerID string, client *cfp.Client) (*ProposalObject, error) {
	submissionStatus := talksv1.ProposalStateDraft
	if obj.Spec.Final {
		submissionStatus = talksv1.ProposalStateFinal
	}

	proposal, err := createProposalPayload(obj, speakerID, submissionStatus)
	if err != nil {
		return nil, err
	}

	// Create the proposal
	resp, err := client.Create(ctx, cfp.ProposalPath, proposal)
	if err != nil {
		return nil, err
	}

	p := &ProposalObject{}
	err = json.Unmarshal(resp, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *ProposalReconciler) updateSubmission(ctx context.Context, obj *talksv1.Proposal, speakerID string, client *cfp.Client) (*ProposalObject, error) {
	switch obj.Status.Submission {
	case talksv1.ProposalStateDraft:
		// If the proposal is marked final, and the submission status is not final, create an entry in cfp.
		if obj.Spec.Final {
			// Create a draft proposal
			proposal, err := createProposalPayload(obj, speakerID, talksv1.ProposalStateFinal)
			if err != nil {
				return nil, err
			}
			// Update the proposal
			resp, err := client.Update(ctx, cfp.ProposalPath, fmt.Sprintf("%s-%s", obj.Namespace, obj.Name), proposal)
			if err != nil {
				return nil, err
			}

			p := &ProposalObject{}
			err = json.Unmarshal(resp, p)
			if err != nil {
				return nil, err
			}
			return p, nil
		} else {
			// Check if the proposal content needs to be updated
			content, err := client.Get(ctx, cfp.ProposalPath, fmt.Sprintf("%s-%s", obj.Namespace, obj.Name))
			if err != nil {
				return nil, err
			}

			// Create a draft proposal
			proposal, err := createProposalPayload(obj, speakerID, talksv1.ProposalStateDraft)
			if err != nil {
				return nil, err
			}

			// If the proposal content is not the same, update the proposal
			if !bytes.Equal(content, proposal) {
				resp, err := client.Update(ctx, cfp.ProposalPath, fmt.Sprintf("%s-%s", obj.Namespace, obj.Name), proposal)
				if err != nil {
					return nil, err
				}
				p := &ProposalObject{}
				err = json.Unmarshal(resp, p)
				if err != nil {
					return nil, err
				}
				return p, nil
			}
		}
	case talksv1.ProposalStateFinal:
		// If the proposal submission is final, no need to update
		// Return
	}
	return nil, nil
}

// reconcileDelete will delete the obj from the CFP API if it is still a draft.
func (r *ProposalReconciler) reconcileDelete(ctx context.Context, obj *talksv1.Proposal, client *cfp.Client) (ctrl.Result, error) {
	// api call to delete the proposal if it is still a draft
	if obj.Status.Submission == talksv1.ProposalStateDraft {
		err := client.Delete(ctx, cfp.ProposalPath, fmt.Sprintf("%s-%s", obj.Namespace, obj.Name))
		if err != nil {
			// return the error so we can requeue
			return ctrl.Result{}, err
		}
	}
	// clean the finalizer
	controllerutil.RemoveFinalizer(obj, talksv1.Finalizer)

	// Stop the reconciliation
	return ctrl.Result{}, nil
}

func createProposalPayload(obj *talksv1.Proposal, speakerID, submission string) ([]byte, error) {
	body := ProposalObject{
		ID:         fmt.Sprintf("%s-%s", obj.Namespace, obj.Name),
		Title:      obj.Spec.Title,
		Abstract:   obj.Spec.Abstract,
		Type:       obj.Spec.Type,
		SpeakerID:  speakerID,
		Final:      obj.Spec.Final,
		Submission: Submission{Status: submission},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshalling body: %w", err)
	}

	return payload, nil
}

// ProposalObject is the object that is sent to/received from the CFP API
type ProposalObject struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Abstract   string     `json:"abstract"`
	Type       string     `json:"type"`
	SpeakerID  string     `json:"speakerID"`
	Final      bool       `json:"final"`
	Submission Submission `json:"submission"`
}

type Submission struct {
	Status     string    `json:"status"`
	LastUpdate time.Time `json:"lastUpdate"`
}
