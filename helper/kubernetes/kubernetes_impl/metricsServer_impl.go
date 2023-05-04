package kubernetes_impl

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"metrics-server-exporter/log"
)

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
