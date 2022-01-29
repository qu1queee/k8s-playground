package controller

import (
	"fmt"
	"time"

	clientset "github.com/shipwright-io/build/pkg/client/clientset/versioned"
	shpscheme "github.com/shipwright-io/build/pkg/client/clientset/versioned/scheme"
	informers "github.com/shipwright-io/build/pkg/client/informers/externalversions/build/v1alpha1"
	listers "github.com/shipwright-io/build/pkg/client/listers/build/v1alpha1"
	"github.com/sirupsen/logrus"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	kubeclientset  kubernetes.Interface
	shpclientset   clientset.Interface
	buildRunLister listers.BuildRunLister
	buildRunSynced cache.InformerSynced
	// use a queue to ensure we process one resource at-a-time, and to avoid
	// multiple workers to work on the same resource
	workqueue workqueue.RateLimitingInterface
}

func NewController(
	kubeclientset kubernetes.Interface,
	shpclientset clientset.Interface,
	buildrunInformer informers.BuildRunInformer,
) *Controller {

	// add Shipwright scheme to the k8s default schemes
	utilruntime.Must(shpscheme.AddToScheme(scheme.Scheme))

	controller := &Controller{
		kubeclientset:  kubeclientset,
		shpclientset:   shpclientset,
		buildRunLister: buildrunInformer.Lister(),
		buildRunSynced: buildrunInformer.Informer().HasSynced,
		workqueue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "BuildRuns"),
	}

	logrus.Infof("Defining BuildRun informer event handler")

	// define the handler for our buildRun informer
	buildrunInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueFoo,
			UpdateFunc: func(old, new interface{}) {
				controller.enqueueFoo(new)
			},
		},
	)
	return controller
}

// take our buildrun object and enqueue it
func (c *Controller) enqueueFoo(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	logrus.Infof("Enqueuing key '%s'", key)
	c.workqueue.Add(key)
}

func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	if ok := cache.WaitForCacheSync(stopCh, c.buildRunSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// read a single item from the queue and then process it by calling the
// syncHandler
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj) // this is relevant, otherwise we might be tracking something forever
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		if err := c.syncHandler(key); err != nil {
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		c.workqueue.Forget(obj)
		logrus.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	logrus.Infof("Namespace: '%s', name: '%s'", namespace, name)

	// Here is the place where the code for processing or reacting on the objects
	// event should be defined.

	return nil
}
