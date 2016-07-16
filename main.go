package main

import (
	"flag"
	"github.com/softputer/kuber-controller/controller"
	"github.com/softputer/kuber-controller/provider"
	"os"
	"os/signal"
	"syscall"	
)

var (
	lbControllerName = flag.String("lb-controller", "kubernetes", "Ingress controller name")
	lbProviderName = flag.String("lb-provider", "haproxy", "Lb Controller name")

	lbc controller.LBController
	lbp controller.LBProvider
)


func setEnv() {
	flag.Parse()
	lbc = controller.GetController(*lbControllerName)
	if lbc == nil {
		logrus.Fatal("Unable to find controller by name %s", *lbControllerName)
	}
	lbp = provider.GetProvider(*lbProviderName)
	if lbp == nil {
		logrus.Fatal("Unable to find provider by name %s", *lbProviderName)	
	}
}

func main() {
	logrus.Infof("Starting Kube LB Service")
	setEnv()
	logrus.Infof("LB controller: %s", lbc.GetName())
	logrus.Infof("LB provider: %s", lbp.GetName())
	
	go handleSigterm(lbc, lbp)

	lbc.Run(lbp)	
}

func handleSigterm(lbc controller.LBController, lbp provider.LBProvider) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	<-signalChan
	logrus.Info("Received SIGTERM, shutting down")

	exitCode := 0
	
	if err := lbc.Stop(); err != nil {
		logrus.Infof("Error during shutdown %v", err)
		exitCode = 1
	}
	logrus.Infof("Exiting with %v", exitCode)
	os.Exit(exitCode)
}
