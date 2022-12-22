package client

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClient(kubeconfig string) (*kubernetes.Clientset, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	return clientSet, err
}
