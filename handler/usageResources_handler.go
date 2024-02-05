package handler

import (
	"fmt"
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
			Help:      "Metrics server pod cpu utilization (m or Millicore)",
		}, []string{"namespace", "kind_owner", "name_owner", "pod", "container", "os"}),
		podsMemory: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "metrics_server",
			Name:      "pod_memory_usage",
			Help:      "Metrics server pod memory utilization (Mi or Mebibyte)",
		}, []string{"namespace", "kind_owner", "name_owner", "pod", "container", "os"}),
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
						log.Error(podMetrics.Name + "," + container.Name + " " + err.Error())
						continue
					}
					memoryMebibytes, err := convertKibibytesToMebibytes(container.Usage.Memory().String())
					if err != nil {
						log.Error(podMetrics.Name + "," + container.Name + " " + err.Error())
						continue
					}

					osType, err := ur.Kubernetes.DetectPodOs(podMetrics.Namespace, podMetrics.Name)
					if err != nil {
						log.Error(err.Error())
					}

					kindOwner, nameOwner, err := ur.Kubernetes.GetPodOwner(podMetrics.Namespace, podMetrics.Name)
					if err != nil {
						log.Error(err.Error())
					}

					m.podsCpu.With(prometheus.Labels{
						"namespace":  podMetrics.Namespace,
						"kind_owner": kindOwner,
						"name_owner": nameOwner,
						"pod":        podMetrics.Name,
						"container":  container.Name,
						"os":         osType}).Set(cpuMillicores)
					m.podsMemory.With(prometheus.Labels{
						"namespace":  podMetrics.Namespace,
						"kind_owner": kindOwner,
						"name_owner": nameOwner,
						"pod":        podMetrics.Name,
						"container":  container.Name,
						"os":         osType}).Set(memoryMebibytes)
				}
			}
			time.Sleep(30 * time.Second)
			//to clear all the previously set metrics before setting new values
			m.podsCpu.Reset()
			m.podsMemory.Reset()
		}
	}()

}

func (ur *UsageResourcesHandler) HandlerMetricsPodUsage() (*v1beta1.PodMetricsList, error) {

	if !strings.Contains(os.Getenv("LOOK_NAMESPACES"), ",") {
		podMetricsList, err := ur.Kubernetes.MetricsServerPodUsage(os.Getenv("LOOK_NAMESPACES"))
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		return podMetricsList, nil
	} else {
		allPodMetricsLists := &v1beta1.PodMetricsList{}
		namespaces := strings.Split(os.Getenv("LOOK_NAMESPACES"), ",")
		for _, namespace := range namespaces {
			podMetricsList, err := ur.Kubernetes.MetricsServerPodUsage(namespace)
			if err != nil {
				log.Error(err.Error())
				return nil, err
			}
			allPodMetricsLists.Items = append(allPodMetricsLists.Items, podMetricsList.Items...)
		}
		return allPodMetricsLists, nil
	}
	return nil, nil
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
			log.Error(err)
			return 0, err
		}
		return float64(cores) / divisorVal, nil
	}

	log.Debug("invalid core string format: " + coreStr)
	return 0, nil
}

// refer ==> https://medium.com/swlh/understanding-kubernetes-resource-cpu-and-memory-units-30284b3cc866
func convertKibibytesToMebibytes(inputStr string) (float64, error) {
	var value int
	var err error

	if strings.HasSuffix(inputStr, "Ki") {
		value, err = strconv.Atoi(strings.TrimSuffix(inputStr, "Ki"))
		if err != nil {
			return 0, err
		}
	} else if strings.HasSuffix(inputStr, "Mi") {
		value, err = strconv.Atoi(strings.TrimSuffix(inputStr, "Mi"))
		if err != nil {
			return 0, err
		}
		value = value * 1024 // Convert mebibytes to kibibytes
	} else {
		return 0, fmt.Errorf("invalid suffix: expected 'Ki' or 'Mi'")
	}

	mebibytes := float64(value) / 1024
	return mebibytes, nil
}
