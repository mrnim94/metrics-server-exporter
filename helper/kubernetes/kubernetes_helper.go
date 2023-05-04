package kubernetes

import "k8s.io/metrics/pkg/apis/metrics/v1beta1"

type Kubernetes interface {
	MetricsServerPodUsage(namespace string) (*v1beta1.PodMetricsList, error)
}
