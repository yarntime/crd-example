package controller

import (
	exampleClient "github.com/yarntime/crd-example/pkg/client/example-client"
	k8s "k8s.io/client-go/kubernetes"
	"time"
)

type Config struct {
	Address               string
	ConcurrentJobHandlers int
	StopCh                chan struct{}
	ResyncPeriod          time.Duration
	K8sClient             *k8s.Clientset
	ExampleClient         *exampleClient.ExampleClient
}
