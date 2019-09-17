package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Define your schema name and the version
var SchemeGroupVersion = schema.GroupVersion{
	Group:   "foghorn",
	Version: "v1alpha1",
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}