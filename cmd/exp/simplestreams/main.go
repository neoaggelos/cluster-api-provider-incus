package main

import (
	"context"
	"os"

	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	ctx        context.Context
	log        = ctrl.Log
	logOptions = logs.NewOptions()
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	ctx = ctrl.SetupSignalHandler()
	ctrl.SetLogger(klog.Background())
	ctx = ctrl.LoggerInto(ctx, log)
}
