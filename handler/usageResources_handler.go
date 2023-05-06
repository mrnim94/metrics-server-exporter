package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"metrics-server-exporter/helper/kubernetes"
	"metrics-server-exporter/log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type UsageResourcesHandler struct {
	Kubernetes kubernetes.Kubernetes
}

type metrics struct {
	podsCpu, podsMemory *prometheus.GaugeVec
}

func (ur *UsageResourcesHandler) NewMetrics(reg prometheus.Registerer) {
	m := &metrics{
		podsCpu: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "metrics_server",
			Name:      "pod_cpu_usage",
			Help:      "Metrics server pod cpu utilization",
		}, []string{"namespace", "pod", "container"}),
		podsMemory: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "metrics_server",
			Name:      "pod_memory_usage",
			Help:      "Metrics server pod memory utilization",
		}, []string{"namespace", "pod", "container"}),
	}
	reg.MustRegister(m.podsCpu)
	reg.MustRegister(m.podsMemory)

	go func() {
		for {
			// set Metrics
			resultPodMetricsList, err := ur.HandlerMetricsPodUsage()
			if err != nil {
				log.Error(err.Error())
				continue
			}
			for _, podMetrics := range resultPodMetricsList.Items {
				for _, container := range podMetrics.Containers {
					cpuMillicores, err := convertNanocoresToMillicores(container.Usage.Cpu().String())
					if err != nil {
						log.Error(err.Error())
						continue
					}
					memoryMebibytes, err := convertKibibytesToMebibytes(container.Usage.Memory().String())
					if err != nil {
						log.Error(err.Error())
						continue
					}
					m.podsCpu.With(prometheus.Labels{"namespace": podMetrics.Namespace, "pod": podMetrics.Name, "container": container.Name}).Set(cpuMillicores)
					m.podsMemory.With(prometheus.Labels{"namespace": podMetrics.Namespace, "pod": podMetrics.Name, "container": container.Name}).Set(memoryMebibytes)
				}
			}
			time.Sleep(30 * time.Second)
		}
	}()

}

func (ur *UsageResourcesHandler) HandlerMetricsPodUsage() (*v1beta1.PodMetricsList, error) {

	podMetricsList, err := ur.Kubernetes.MetricsServerPodUsage(os.Getenv("LOOK_NAMESPACE"))
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return podMetricsList, nil
}

// HandlerPodUsage /*** call API
func (ur *UsageResourcesHandler) HandlerPodUsage(c echo.Context) error {

	podMetricsList, err := ur.Kubernetes.MetricsServerPodUsage("LOOK_NAMESPACE")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return c.JSON(http.StatusOK, podMetricsList.Items)
}

func convertNanocoresToMillicores(coreStr string) (float64, error) {
	divisor := map[string]float64{
		"n": 1000000.0,
		"m": 1.0,
	}

	suffix := strings.TrimLeftFunc(coreStr, unicode.IsDigit)
	coreStr = strings.TrimSuffix(coreStr, suffix)

	if divisorVal, ok := divisor[suffix]; ok {
		cores, err := strconv.ParseInt(coreStr, 10, 64)
		if err != nil {
			return 0, err
		}
		return float64(cores) / divisorVal, nil
	}

	iErr := errors.New("invalid core string format: " + coreStr)
	log.Error(iErr)
	return 0, iErr
}

// refer ==> https://medium.com/swlh/understanding-kubernetes-resource-cpu-and-memory-units-30284b3cc866
func convertKibibytesToMebibytes(kibibytesStr string) (float64, error) {
	kibibytes, err := strconv.Atoi(strings.TrimSuffix(kibibytesStr, "Ki"))
	if err != nil {
		return 0, err
	}
	mebibytes := float64(kibibytes) / 1024
	return mebibytes, nil
}
