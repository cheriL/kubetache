package kube

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type Client interface {
	Run()
	Stop()

	GetPod(string, string) *corev1.Pod
	GetReplicaSet(string, string) *appsv1.ReplicaSet
	GetServices(map[string]string, string) []*corev1.Service
}
