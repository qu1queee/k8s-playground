package controller

import (
	clientset "github.com/shipwright-io/build/pkg/client/clientset/versioned"
	informers "github.com/shipwright-io/build/pkg/client/informers/externalversions"
	"k8s.io/client-go/kubernetes"
)

// Controller is the controller implementation for Foo resources
type Controller struct {
}

func NewController(
	kubeclientset kubernetes.Interface,
	shpclientset clientset.Interface,
	shpInformer informers.SharedInformerFactory,
) *Controller {
	return nil
}

func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	return nil
}
