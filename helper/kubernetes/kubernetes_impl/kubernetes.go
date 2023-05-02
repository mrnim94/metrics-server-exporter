package kubernetes_impl

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"metrics-server-exporter/helper/kubernetes"
	"metrics-server-exporter/log"
)

type KubeConfiguration struct {
	KubeConfig string
}

func NewKubernetesConnection(kc *KubeConfiguration) kubernetes.Kubernetes {
	return &KubeConfiguration{
		KubeConfig: kc.KubeConfig,
	}
}

func (kc *KubeConfiguration) accessKubernetes() (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kc.KubeConfig)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return config, nil
}
