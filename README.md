# Metrics Server Exporter

The metrics server exporter is solution to help you monitor RAM/CPU easily.  
You can monitor RAM/CPU of the windows container via the metrics server exporter


## Prequirement

To the metrics-server-exporter, you need prepare something:

 - Install [Kubernetes Metrics Server](https://github.com/kubernetes-sigs/metrics-server)

## Install metrics-server-exporter   

We will install the metrics-server-exporter through Helm Chart.

```
helm repo add openvscode-server-helm https://mrnim94.github.io/metrics-server-exporter

helm search repo metrics-server-exporter
NAME                                            CHART VERSION   APP VERSION     DESCRIPTION
metrics-server-exporter/metrics-server-exporter 0.1.2           0.0.1           A Helm chart for Kubernetes
```
### Importance
Currently, the metrics-server-exporter won't collect all pods metrics in your k8s. It will only collect metrics following the namespace that you specify.  

#### VALUES FILES
```
#values.yaml
envVars:
  LOOK_NAMESPACE: <namespace>
```

### Metrics
```
http://localhost:1994
```
Name | Description | Type | Labels
-----|-------------|------|-------
`metrics_server_pod_cpu_usage` | Metrics server pod cpu utilization (m or Millicore) | gauge | `container`
`metrics_server_pod_memory_usage` | Metrics server pod memory utilization (Mi or Mebibyte) | gauge | `container`

### Grafana Daskboard
- Get [Grafana Json](https://raw.githubusercontent.com/mrnim94/metrics-server-exporter/master/grafana/Pods%20-%20Metrics%20Server%20Monitor%20Prometheus-1688617368230.json)


### Refer to:  

 - [Pods - Metrics Server Monitor Prometheus](https://grafana.com/grafana/dashboards/8760-pods-metrics-server-monitor-prometheus/)
 - [metrics-server-monitor-prometheus](https://github.com/AdrianBalcan/metrics-server-monitor-prometheus)

#### Create Helm Package

```sh
helm package ./helm-chart/metrics-server-exporter --destination ./helm-chart/
helm repo index . --url https://mrnim94.github.io/metrics-server-exporter
```
