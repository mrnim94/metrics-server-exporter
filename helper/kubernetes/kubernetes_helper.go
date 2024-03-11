package kubernetes

import (
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"metrics-server-exporter/model"
)

type Kubernetes interface {
	MetricsServerPodUsage(namespace string) (*v1beta1.PodMetricsList, error)
	DetectPodOs(namespace, podName string) (string, error)
	GetPodOwner(namespace, podName string) (ownerKind, ownerName string, err error)
	CurrentMetricsForHPAs(namespace string) (model.HPAUtilizations, error)
}
