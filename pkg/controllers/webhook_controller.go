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
	apierrs "k8s.io/apimachinery/pkg/api/errors"

	foghornv1alpha1 "github.com/foghornci/foghorn/pkg/apis/foghorn/v1alpha1"
)

// WebhookReconciler reconciles a Webhook object
type WebhookReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=foghorn,resources=webhooks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=foghorn,resources=webhooks/status,verbs=get;update;patch

func (r *WebhookReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("webhook", req.NamespacedName)

	var webhook foghornv1alpha1.Webhook
	if err := r.Get(ctx, req.NamespacedName, &webhook); err != nil {
		return ctrl.Result{}, ignoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

func (r *WebhookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&foghornv1alpha1.Webhook{}).
		Complete(r)
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
			return nil
	}
	return err
}
