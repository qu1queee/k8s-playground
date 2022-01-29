package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qu1queee/k8s-playground/cc/pkg/controller"
	clientset "github.com/shipwright-io/build/pkg/client/clientset/versioned"
	informers "github.com/shipwright-io/build/pkg/client/informers/externalversions"
	"github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var onlyOneSignalHandler = make(chan struct{})
var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

func main() {
	logrus.Infof("Starting CC controller")

	// this calls a func that the controller-runtime already implements
	// with a context, see https://github.com/kubernetes-sigs/controller-runtime/blob/master/pkg/manager/signals/signal.go#L27-L45
	stopCh := SetupSignalHandler()

	logrus.Infof("Getting k8s config")
	cfg, err := config.GetConfig()
	if err != nil {
		os.Exit(1)
	}

	kubeClient, _ := kubernetes.NewForConfig(cfg)
	shpclient, _ := clientset.NewForConfig(cfg)

	exampleInformerFactory := informers.NewSharedInformerFactory(shpclient, time.Second*30)

	controller := controller.NewController(
		kubeClient,
		shpclient,
		exampleInformerFactory.Shipwright().V1alpha1().BuildRuns(),
	)

	// run the informers
	exampleInformerFactory.Start(stopCh)

	// run the controller with two workers, basically two separate
	// go routines that will consume events from the queue
	if err = controller.Run(2, stopCh); err != nil {
		logrus.Error("Error running controller: %s", err.Error())
	}

}

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler)

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1)
	}()

	return stop
}
