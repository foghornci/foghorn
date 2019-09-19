/*

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

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	foghornv1alpha1 "github.com/foghornci/foghorn/pkg/apis/foghorn.jenkins.io/v1alpha1"
)

// GitEventReconciler reconciles a GitEvent object
type GitEventReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=foghorn.jenkins.io,resources=gitevents,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=foghorn.jenkins.io,resources=gitevents/status,verbs=get;update;patch

func (r *GitEventReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("gitevent", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *GitEventReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&foghornv1alpha1.GitEvent{}).
		Complete(r)
}
