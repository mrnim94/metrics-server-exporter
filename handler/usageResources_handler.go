package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"metrics-server-exporter/helper/kubernetes"
	"metrics-server-exporter/log"
	"net/http"
	"strconv"
	"strings"
)

type UsageResourcesHandler struct {
	Echo       *echo.Context
	Kubernetes kubernetes.Kubernetes
}

type metrics struct {
	podsCpu *prometheus.GaugeVec
}

func (ur *UsageResourcesHandler) NewMetrics(reg prometheus.Registerer) {
	m := &metrics{
		podsCpu: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "metricsServer",
			Name:      "metrics_server_pod_cpu_usage",
			Help:      "Information about the My App environment.",
		},
			[]string{"namespace", "pod", "container"}),
	}
	reg.MustRegister(m.podsCpu)

	// set Metrics
	resultPodMetricsList, err := ur.HandlerMetricsPodUsage()
	if err != nil {
		log.Error(err.Error())
	}
	for _, podMetrics := range resultPodMetricsList.Items {
		for _, container := range podMetrics.Containers {
			cpuMillicores, err := convertNanocoresToMillicores(container.Usage.Cpu().String())
			if err != nil {
				log.Error(err.Error())
				return
			}
			m.podsCpu.With(prometheus.Labels{"namespace": podMetrics.Namespace, "pod": podMetrics.Name, "container": container.Name}).Set(cpuMillicores)
		}
	}
}

func (ur *UsageResourcesHandler) HandlerMetricsPodUsage() (*v1beta1.PodMetricsList, error) {

	podMetricsList, err := ur.Kubernetes.MetricsServerPodCpuUsage("openvscode-server")
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return podMetricsList, nil
}

// HandlerPodUsage /*** call API
func (ur *UsageResourcesHandler) HandlerPodUsage(c echo.Context) error {

	podMetricsList, err := ur.Kubernetes.MetricsServerPodCpuUsage("openvscode-server")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return c.JSON(http.StatusOK, podMetricsList.Items)
}

func convertNanocoresToMillicores(nanocoresStr string) (float64, error) {
	nanocoresStr = strings.TrimSuffix(nanocoresStr, "n")
	nanocores, err := strconv.ParseInt(nanocoresStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return float64(nanocores) / 1000000.0, nil
}
