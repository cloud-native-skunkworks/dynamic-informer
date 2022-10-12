package main

import (
	"log"
	"os"
	"time"

	"github.com/open-feature/open-feature-operator/apis/core/v1alpha1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func buildConfiguration() (*rest.Config, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	var clusterConfig *rest.Config
	var err error
	if kubeconfig != "" {
		clusterConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		clusterConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	return clusterConfig, nil
}

func main() {
	clusterConfig, err := buildConfiguration()
	if err != nil {
		panic(err)
	}

	clusterClient, err := dynamic.NewForConfig(clusterConfig)
	if err != nil {
		log.Fatalln(err)
	}

	resource := v1alpha1.GroupVersion.WithResource("featureflagconfigurations")
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clusterClient,
		time.Minute, "", nil)
	informer := factory.ForResource(resource).Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Println("AddFunc")
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			log.Println("UpdateFunc")
		},
		DeleteFunc: func(obj interface{}) {
			log.Println("DeleteFunc")
		},
	})

	informer.Run(make(chan struct{}))
}
