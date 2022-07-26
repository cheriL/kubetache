package kube

import (
	corev1 "k8s.io/api/core/v1"
)

type Client interface {
	Run()
	Stop()

	GetPod(name, namespace string) *corev1.Pod
}
