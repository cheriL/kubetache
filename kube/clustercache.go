package kube

import (
	appsv1 "k8s.io/api/apps/v1"
	autoscalingV1 "k8s.io/api/autoscaling/v1"
	autoscalingV2 "k8s.io/api/autoscaling/v2beta1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
)

type ClusterCache struct {
	k      *kubernetes.Clientset
	stopCh chan struct{}

	pod          CacheController
	service      CacheController
	deployment   CacheController
	replicaset   CacheController
	statefulSet  CacheController
	daemonSet    CacheController
	job          CacheController
	cronJob      CacheController
	configmap    CacheController
	pv           CacheController
	pvc          CacheController
	hpaV1        CacheController
	hpaV2Beta2   CacheController
	//storageClass CacheController
}

func NewClusterCache(client *kubernetes.Clientset) Client {
	coreClient  := client.CoreV1().RESTClient()
	appClient   := client.AppsV1().RESTClient()
	batchClient := client.BatchV1().RESTClient()

	hpaV2Beta2Client := client.AutoscalingV2beta2().RESTClient()
	hpaV1Client      := client.AutoscalingV1().RESTClient()

	return &ClusterCache{
		k:           client,
		pod:         NewController(coreClient, &corev1.Pod{}, "pods", ""),
		service:     NewController(coreClient, &corev1.Service{}, "services", ""),
		deployment:  NewController(appClient, &appsv1.Deployment{}, "deployments", ""),
		replicaset:  NewController(appClient, &appsv1.ReplicaSet{}, "replicasets", ""),
		statefulSet: NewController(appClient, &appsv1.StatefulSet{}, "statefulsets", ""),
		daemonSet:   NewController(appClient, &appsv1.DaemonSet{}, "daemonsets", ""),
		job:         NewController(batchClient, &batchv1.Job{}, "jobs", ""),
		cronJob:     NewController(batchClient, &batchv1.CronJob{}, "cronjobs", ""),
		configmap:   NewController(coreClient, &corev1.ConfigMap{}, "configmaps", ""),
		pv:          NewController(coreClient, &corev1.PersistentVolume{}, "persistentvolumes", ""),
		pvc:         NewController(coreClient, &corev1.PersistentVolumeClaim{}, "persistentvolumeclaims", ""),
		hpaV1:       NewController(hpaV1Client, &autoscalingV1.HorizontalPodAutoscaler{}, "horizontalpodautoscalers", ""),
		hpaV2Beta2:  NewController(hpaV2Beta2Client, &autoscalingV2.HorizontalPodAutoscaler{}, "horizontalpodautoscalers", ""),
	}
}

func (cc *ClusterCache)Run() {
	if cc.stopCh != nil {
		return
	}
	stopCh := make(chan struct{})

	go cc.pod.Run(stopCh)
	go cc.replicaset.Run(stopCh)
	go cc.service.Run(stopCh)

	cc.stopCh = stopCh
}

func (cc *ClusterCache)Stop() {
	if cc.stopCh == nil {
		return
	}

	close(cc.stopCh)
	cc.stopCh = nil
}

func (cc *ClusterCache)GetPod(name, namespace string) *corev1.Pod {

	podName, ns := name, namespace
	if podName == "" {
		return nil
	}

	if namespace == "" {
		ns = corev1.NamespaceDefault
	}

	objList := cc.pod.ListObjs()
	for _, obj := range objList {
		if pod, ok := obj.(*corev1.Pod); ok {
			if pod.Name == podName && pod.Namespace == ns {
				return pod
			}
		}
	}

	return nil
}

func (cc *ClusterCache)GetReplicaSet(name, namespace string) *appsv1.ReplicaSet {

	replicaSet, ns := name, namespace

	if replicaSet == "" {
		return nil
	}

	if namespace == "" {
		ns = corev1.NamespaceDefault
	}

	objList := cc.replicaset.ListObjs()
	for _, obj := range objList {
		if r, ok := obj.(*appsv1.ReplicaSet); ok {
			if r.Name == replicaSet && r.Namespace == ns {
				return r
			}
		}
	}

	return nil
}

func (cc *ClusterCache)GetServices(labels map[string]string, namespace string) []*corev1.Service {
	ns := namespace
	if ns == "" {
		ns = corev1.NamespaceDefault
	}

	var services []*corev1.Service
	labelSet := fields.Set(labels)

	objList := cc.service.ListObjs()
	for _, obj := range objList {
		if s, ok := obj.(*corev1.Service); ok {
			if s.Namespace == ns && s.Spec.Selector != nil {
				selector := fields.SelectorFromSet(s.Spec.Selector)
				if match := selector.Matches(labelSet); match {
					services = append(services, s)
				}
			}
		}
	}

	return services
}