package kubernetes_impl

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"metrics-server-exporter/helper/kubernetes"
	"metrics-server-exporter/log"
	"os"
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
	var config *rest.Config
	_, err := os.Stat(kc.KubeConfig)
	if os.IsNotExist(err) {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Error("Failed to create in-cluster configuration: ", err)
			return nil, err
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kc.KubeConfig)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
	}
	return config, nil
}
