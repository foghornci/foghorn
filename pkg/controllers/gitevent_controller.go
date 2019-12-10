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
	"fmt"

	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	foghornv1alpha1 "github.com/foghornci/foghorn/pkg/apis/foghorn.jenkins.io/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GitEventReconciler reconciles a GitEvent object
type GitEventReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=foghorn.jenkins.io,resources=gitevents,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=foghorn.jenkins.io,resources=gitevents/status,verbs=get;update;patch

// Reconcile brings the actual state of the world to our declared state
func (r *GitEventReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("gitevent", req.NamespacedName)

	// your logic here
	logrus.Infof("reconciling the following resource: %s/%s", req.Namespace, req.Name)

	actions := &foghornv1alpha1.ActionList{}
	err := r.List(ctx, actions)
	if err != nil {
		logrus.WithError(err).Warn("failed to list Action resources")
	}

	// List Action objs to see if one already has this GitEvent, return if so
	for _, action := range actions.Items {
		name := fmt.Sprintf("%s/%s", action.Spec.ParentEvent.Namespace, action.Spec.ParentEvent.Name)
		reqName := fmt.Sprintf("%s/%s", req.Namespace, req.Name)
		if name == reqName {
			return ctrl.Result{}, nil
		}
	}

	// Get the GitEvent and assign it to our new Action
	gitEvent := &foghornv1alpha1.GitEvent{}
	err = r.Get(ctx, req.NamespacedName, gitEvent)
	if err != nil {
		logrus.WithError(err).Warn("failed to fetch GitEvent Resource")
	}

	a := &foghornv1alpha1.Action{
		Spec: foghornv1alpha1.ActionSpec{
			ParentEvent: gitEvent,
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "action-",
			Namespace:    req.Namespace,
		},
	}

	err = r.Create(ctx, a)
	if err != nil {
		logrus.WithError(err).Warn("failed to create Action resource")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller
func (r *GitEventReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&foghornv1alpha1.GitEvent{}).
		Complete(r)
}
