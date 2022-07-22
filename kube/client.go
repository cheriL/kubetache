package kube

import (
	v1 "k8s.io/api/core/v1"
	rt "k8s.io/apimachinery/pkg/runtime"
)

type Client interface {
	CheckResource(resource, namespace string, resourceType rt.Object) bool
	
	GetPod(name, namespace string) *v1.Pod
}
