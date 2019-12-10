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
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/jenkins-x/go-scm/scm"
	scmfactory "github.com/jenkins-x/go-scm/scm/factory"
	"github.com/sirupsen/logrus"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	foghornv1alpha1 "github.com/foghornci/foghorn/pkg/apis/foghorn.jenkins.io/v1alpha1"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// ActionReconciler reconciles a Action object
type ActionReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=foghorn.jenkins.io,resources=actions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=foghorn.jenkins.io,resources=actions/status,verbs=get;update;patch

func (r *ActionReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("action", req.NamespacedName)

	// your logic here

	action := &foghornv1alpha1.Action{}
	err := r.Get(ctx, req.NamespacedName, action)
	if err != nil {
		logrus.WithError(err).Warn("failed to fetch Action resource %s/%s", req.Namespace, req.Name)
	}

	repo := action.Spec.ParentEvent.Spec.ParsedWebhook.Webhook.Repository()
	configMapNameSlice := strings.Split(repo.FullName, "/")
	if len(configMapNameSlice) < 2 {
		return ctrl.Result{}, fmt.Errorf("error parsing repo name")
	}
	configMapName := fmt.Sprintf("%s-%s", configMapNameSlice[0], configMapNameSlice[1])
	cm := &corev1.ConfigMap{}
	namespacedName := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      configMapName,
	}
	err = r.Get(ctx, namespacedName, cm)
	if err != nil {
		logrus.WithError(err).Warnf("failed to fetch ConfigMap %s/%s", req.Namespace, configMapName)
	}
	if err := r.createConfigMapFromRepoConfig(req.Namespace, configMapName, repo); err != nil {
		logrus.WithError(err).Warnf("failed to create ConfigMap %s/%s", req.Namespace, configMapName)
	}
	return ctrl.Result{}, nil
}

func (r *ActionReconciler) createConfigMapFromRepoConfig(namespace, name string, repo scm.Repository) error {
	ctx := context.Background()
	oauthToken := os.Getenv("OAUTH_TOKEN")
	if oauthToken == "" {
		return fmt.Errorf("no git oauth present")
	}

	gitKind := os.Getenv("GIT_KIND")
	if gitKind == "" {
		gitKind = "github"
	}

	gitBaseURL := os.Getenv("GIT_BASE_URL")

	scmClient, err := scmfactory.NewClient(gitKind, gitBaseURL, oauthToken)
	if err != nil {
		return err
	}

	content, _, err := scmClient.Contents.Find(ctx, repo.FullName, "lighthouse.yaml", "fork")
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}

	data := map[string]string{}
	err = yaml.Unmarshal(content.Data, data)
	if err != nil {
		return err
	}

	cm.Data = data

	err = r.Create(context.Background(), cm)
	if err != nil {
		return err
	}
	return nil
}

func (r *ActionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&foghornv1alpha1.Action{}).
		Complete(r)
}
