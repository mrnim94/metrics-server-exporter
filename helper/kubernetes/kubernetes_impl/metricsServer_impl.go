package kubernetes_impl

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

func (kc *KubeConfiguration) MetricsServerPodCpuUsage(namespace string) (*v1beta1.PodMetricsList, error) {

	config, err := kc.accessKubernetes()
	metricsClientset, err := metricsv.NewForConfig(config)
	if err != nil {
		fmt.Printf("error %s, getting inclusterconfig", err.Error())
	}
	podMetricsList, err := metricsClientset.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return podMetricsList, nil
}
