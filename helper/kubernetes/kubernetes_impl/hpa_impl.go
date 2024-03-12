package kubernetes_impl

import (
	"context"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"metrics-server-exporter/log"
	"metrics-server-exporter/model"
)

func (kc *KubeConfiguration) CurrentMetricsForHPAs(namespace string) (model.HPAUtilizations, error) { //NOSONAR - Refactor this method to reduce its Cognitive Complexity

	results := model.HPAUtilizations{}

	// Load kubeconfig
	config, err := kc.accessKubernetes()
	if err != nil {
		log.Error(err)
		return results, err
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err)
		return results, err
	}
	hpaList, err := clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Error("Error fetching HPAs: %v", err)
		return results, err
	}
	for _, hpa := range hpaList.Items {
		log.Debugf("HPA: %s", hpa.Name)
		log.Debugf("ScaleTargetRef: Kind=%s, Name=%s", hpa.Spec.ScaleTargetRef.Kind, hpa.Spec.ScaleTargetRef.Name)

		// Display current metrics for this HPA
		for _, metric := range hpa.Status.CurrentMetrics {
			switch metric.Type {
			case autoscalingv2.ResourceMetricSourceType:
				log.Debugf("  Resource Metric: %s, Current Average Utilization: %d, Current Average Value: %s",
					metric.Resource.Name, *metric.Resource.Current.AverageUtilization, metric.Resource.Current.AverageValue.String())
				results = append(results, model.HPAUtilization{
					MetricName:         string(metric.Resource.Name),
					MetricType:         string(metric.Type),
					HPAOwner:           hpa.Name,
					ScaleTargetRefKind: hpa.Spec.ScaleTargetRef.Kind,
					ScaleTargetRefName: hpa.Spec.ScaleTargetRef.Name,
					AverageUtilization: float64(*metric.Resource.Current.AverageUtilization),
				})

			case autoscalingv2.ExternalMetricSourceType:
				var correctNum int
				for iMetricOfSpec, metricOfSpec := range hpa.Spec.Metrics {
					if metricOfSpec.External.Metric.Name == metric.External.Metric.Name {
						correctNum = iMetricOfSpec
						break
					}
				}
				log.Debugf("  External Metric: %s, Define Average Value: %s, Current Average Value: %s",
					metric.External.Metric.Name, hpa.Spec.Metrics[correctNum].External.Target.AverageValue.String(), metric.External.Current.AverageValue.String())

				// Convert to float64 using AsApproximateFloat64
				currentAverageValueFloat := metric.External.Current.AverageValue.AsApproximateFloat64()
				targetAverageValueFloat := hpa.Spec.Metrics[correctNum].External.Target.AverageValue.AsApproximateFloat64()

				// Now you can perform calculations with these float64 values
				if targetAverageValueFloat != 0 {
					utilizationTemp := currentAverageValueFloat / targetAverageValueFloat
					// Convert utilizationTemp to a percentage or format as needed
					utilizationPercentage := utilizationTemp * 100 // Example to convert to percentage

					// Now you can use utilizationPercentage safely
					results = append(results, model.HPAUtilization{
						MetricName:         metric.External.Metric.Name,
						MetricType:         string(metric.Type),
						HPAOwner:           hpa.Name,
						ScaleTargetRefKind: hpa.Spec.ScaleTargetRef.Kind,
						ScaleTargetRefName: hpa.Spec.ScaleTargetRef.Name,
						AverageUtilization: utilizationPercentage,
					})
				} else {
					// Handle the error scenario, maybe log that the target value is zero
					log.Error("Target average value is zero")
				}
			}
		}
	}
	return results, nil
}
