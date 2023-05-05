# Metrics Server Exporter

The metrics server exporter is solution to help you monitor RAM/CPU easily.  
You can monitor RAM/CPU of the windows container via the metrics server exporter




## Prequirement

To the metrics-server-exporter, you need prepare something:

 - Install [Kubernetes Metrics Server](https://github.com/kubernetes-sigs/metrics-server)

### Refer to:  

 - [Pods - Metrics Server Monitor Prometheus](https://grafana.com/grafana/dashboards/8760-pods-metrics-server-monitor-prometheus/)
 - [metrics-server-monitor-prometheus](https://github.com/AdrianBalcan/metrics-server-monitor-prometheus)

#### Create Helm Package

```sh
helm package ./helm-chart/* --destination ./helm-chart/
helm repo index . --url https://mrnim94.github.io/metrics-server-exporter/helm-chart
```
