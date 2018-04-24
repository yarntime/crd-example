package main

import (
	"flag"
	"github.com/golang/glog"
	exampleClient "github.com/yarntime/crd-example/pkg/client/example-client"
	k8sclient "github.com/yarntime/crd-example/pkg/client/k8s_client"
	c "github.com/yarntime/crd-example/pkg/controller"
	"github.com/yarntime/crd-example/pkg/tools"
	"time"
)

var (
	apiserverAddress      string
	concurrentJobHandlers int
	resyncPeriod          time.Duration
	baseImage             string
	jobNamespace          string
)

func init() {
	flag.StringVar(&apiserverAddress, "apiserver_address", "", "Kubernetes apiserver address")
	flag.IntVar(&concurrentJobHandlers, "concurrent_job_handlers", 4, "Concurrent job handlers")
	flag.DurationVar(&resyncPeriod, "resync_period", time.Minute*30, "Resync period")
	flag.Set("alsologtostderr", "true")
	flag.Parse()
}

func main() {
	stop := make(chan struct{})

	restConfig, err := tools.GetClientConfig(apiserverAddress)
	if err != nil {
		panic(err.Error())
	}

	glog.Info("register example.")
	err = exampleClient.RegisterExample(restConfig)
	if err != nil {
		panic(err.Error())
	}

	config := &c.Config{
		Address:               apiserverAddress,
		ConcurrentJobHandlers: concurrentJobHandlers,
		ResyncPeriod:          resyncPeriod,
		StopCh:                stop,
		K8sClient:             k8sclient.NewK8sClint(restConfig),
		ExampleClient:         exampleClient.NewExampleClient(restConfig),
	}

	mtc := c.NewExampleController(config)

	glog.Info("run controller.")
	go mtc.Run(stop)

	select {}
}
