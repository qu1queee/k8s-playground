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

func main() {
	logrus.Infof("Starting CC controller")

	stopCh := SetupSignalHandler()

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
		exampleInformerFactory,
	)

	exampleInformerFactory.Start(stopCh)
	if err = controller.Run(2, stopCh); err != nil {
		logrus.Error("Error running controller: %s", err.Error())
	}

}

var onlyOneSignalHandler = make(chan struct{})
var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

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
