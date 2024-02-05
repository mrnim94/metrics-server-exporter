package kubernetes_impl

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"metrics-server-exporter/log"
)

/*
**
This function based on:
> kubectl top pod -n mdaas-engines-prod
NAME                                             CPU(cores)   MEMORY(bytes)
doctor-doom-fb5fd5f94-mwc6g                      0m           166Mi
eks-sample-windows-deployment-5c8c7d6c9c-4frkp   0m           104Mi
eng-7z-7766964875-4mplq                          13m          226Mi
*/
func (kc *KubeConfiguration) MetricsServerPodUsage(namespace string) (*v1beta1.PodMetricsList, error) {

	config, err := kc.accessKubernetes()
	metricsClientset, err := metricsv.NewForConfig(config)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	podMetricsList, err := metricsClientset.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return podMetricsList, nil
}
