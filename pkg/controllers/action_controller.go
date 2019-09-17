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

	foghornv1alpha1 "github.com/foghornci/foghorn/pkg/apis/foghorn/v1alpha1"
)

// ActionReconciler reconciles a Action object
type ActionReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=foghorn,resources=actions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=foghorn,resources=actions/status,verbs=get;update;patch

func (r *ActionReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("action", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *ActionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&foghornv1alpha1.Action{}).
		Complete(r)
}
