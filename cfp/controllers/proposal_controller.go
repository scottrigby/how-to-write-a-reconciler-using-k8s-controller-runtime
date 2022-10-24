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
	"net/http"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	talksv1 "github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime/cfp/api/v1"
)

// ProposalReconciler reconciles a Proposal object
type ProposalReconciler struct {
	client.Client
	HTTPClient     *http.Client
	ControllerName string
	CfpAPI         string
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProposalReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&talksv1.Proposal{}).
		Complete(r)
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
	return
}
