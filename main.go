package main

import (
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"metrics-server-exporter/handler"
	"metrics-server-exporter/helper/kubernetes/kubernetes_impl"
	"metrics-server-exporter/log"
	"metrics-server-exporter/router"
	"os"
)

func init() {
	os.Setenv("APP_NAME", "XXXX")
	log.InitLogger(false)
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
}

func main() {
	kubeconfig := flag.String("kubeconfig", ".kube/config", "location to your confighihi file")
	kube := &kubernetes_impl.KubeConfiguration{KubeConfig: *kubeconfig}

	e := echo.New()

	usageResourcesHandler := handler.UsageResourcesHandler{
		Kubernetes: kubernetes_impl.NewKubernetesConnection(kube),
	}

	// Prometheus client
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector())
	usageResourcesHandler.NewMetrics(reg)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})

	api := router.API{
		Echo:                  e,
		PromHandler:           promHandler,
		UsageResourcesHandler: usageResourcesHandler,
	}
	api.SetupRouter()
	e.Logger.Fatal(e.Start(":1323"))
}
