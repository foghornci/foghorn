/*
Copyright The Kubernetes Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/foghornci/foghorn/pkg/apis/foghorn.jenkins.io/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// GitEventLister helps list GitEvents.
type GitEventLister interface {
	// List lists all GitEvents in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.GitEvent, err error)
	// GitEvents returns an object that can list and get GitEvents.
	GitEvents(namespace string) GitEventNamespaceLister
	GitEventListerExpansion
}

// gitEventLister implements the GitEventLister interface.
type gitEventLister struct {
	indexer cache.Indexer
}

// NewGitEventLister returns a new GitEventLister.
func NewGitEventLister(indexer cache.Indexer) GitEventLister {
	return &gitEventLister{indexer: indexer}
}

// List lists all GitEvents in the indexer.
func (s *gitEventLister) List(selector labels.Selector) (ret []*v1alpha1.GitEvent, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GitEvent))
	})
	return ret, err
}

// GitEvents returns an object that can list and get GitEvents.
func (s *gitEventLister) GitEvents(namespace string) GitEventNamespaceLister {
	return gitEventNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// GitEventNamespaceLister helps list and get GitEvents.
type GitEventNamespaceLister interface {
	// List lists all GitEvents in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.GitEvent, err error)
	// Get retrieves the GitEvent from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.GitEvent, error)
	GitEventNamespaceListerExpansion
}

// gitEventNamespaceLister implements the GitEventNamespaceLister
// interface.
type gitEventNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all GitEvents in the indexer for a given namespace.
func (s gitEventNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.GitEvent, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GitEvent))
	})
	return ret, err
}

// Get retrieves the GitEvent from the indexer for a given namespace and name.
func (s gitEventNamespaceLister) Get(name string) (*v1alpha1.GitEvent, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("gitevent"), name)
	}
	return obj.(*v1alpha1.GitEvent), nil
}
