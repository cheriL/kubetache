package kube

import (
	"fmt"
	"k8s.io/apimachinery/pkg/fields"
	rt "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type CacheController interface {
	Run(chan struct{})

	ListObjs() []interface{}
}

type Controller struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller

	resource string
}

func NewController(client rest.Interface, resourceType rt.Object, resource string, namespace string) CacheController {
	// create the resource watcher
	resListWatcher := cache.NewListWatchFromClient(client, resource, namespace, fields.Everything())

	// create the workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the resource key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the resource than the version which was responsible for triggering the update.
	indexer, informer := cache.NewIndexerInformer(resListWatcher, resourceType, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}, cache.Indexers{})

	return &Controller{
		informer: informer,
		indexer:  indexer,
		queue:    queue,
		resource: resource,
	}
}

func (c *Controller) Run(stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer c.queue.ShutDown()
	klog.Info("starting controller: ", c.resource)

	go c.informer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	//for i := 0; i < threadiness; i++ {
	//	go wait.Until(c.runWorker, time.Second, stopCh)
	//}

	<-stopCh
	klog.Info("stopping controller: ", c.resource)
}

func (c *Controller) ListObjs() []interface{} {
	list := c.indexer.List()

	newList := make([]interface{}, 0, len(list))

	for _, v := range list {
		if obj, ok := v.(rt.Object); ok {
			newList = append(newList, obj.DeepCopyObject())
		}
	}

	return newList
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {
	return true
}
