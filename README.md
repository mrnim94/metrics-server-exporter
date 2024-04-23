## Metrics Server Exporter

The metrics server exporter is solution to help you monitor RAM/CPU easily.  
You can monitor RAM/CPU of the windows container via the metrics server exporter

## Badge

Helm Chart:  
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/metrics-server-exporter)](https://artifacthub.io/packages/search?repo=metrics-server-exporter)

Docker:  
![Docker Image Version](https://img.shields.io/docker/v/mrnim94/metrics-server-exporter)  
![Docker Pulls](https://img.shields.io/docker/pulls/mrnim94/metrics-server-exporter)

[![Quality gate](https://sonarcloud.io/api/project_badges/quality_gate?project=mrnim94_metrics-server-exporter)](https://sonarcloud.io/summary/new_code?id=mrnim94_metrics-server-exporter)  
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmrnim94%2Fmetrics-server-exporter.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmrnim94%2Fmetrics-server-exporter?ref=badge_large&issueType=license)

## Prequirement

To the metrics-server-exporter, you need prepare something:

*   Install [Kubernetes Metrics Server](https://github.com/kubernetes-sigs/metrics-server)

## Install metrics-server-exporter

We will install the metrics-server-exporter through Helm Chart.

```plaintext
helm repo add openvscode-server-helm https://mrnim94.github.io/metrics-server-exporter

helm search repo metrics-server-exporter
NAME                                            CHART VERSION   APP VERSION     DESCRIPTION
metrics-server-exporter/metrics-server-exporter 0.1.2           0.0.1           A Helm chart for Kubernetes
```

### Importance

Default `metrics-server-exporter` will collect the usage resources of ALL namespace

Besides, you can specify many namespace via environment and `metrics-server-exporter` will collect the metrics follow  
your mean

#### VALUES FILES

```plaintext
#values.yaml
envVars:
  LOOK_NAMESPACES: <namespace1>,<namespace2>
```

### Metrics

```plaintext
http://localhost:1994
```

| Name | Description | Type | Labels |
| --- | --- | --- | --- |
| `metrics_server_pod_cpu_usage` | Metrics server pod cpu utilization (m or Millicore) | gauge | `namespace`, `kind_owner`, `name_owner`, `pod`, `container`, `os` |
| `metrics_server_pod_memory_usage` | Metrics server pod memory utilization (Mi or Mebibyte) | gauge | `namespace`, `kind_owner`, `name_owner`, `pod`, `container`, `os` |
| `metrics_server_hpa_utilization` | Current Average Utilization Percentage of each metric that is created by HPA (%) | gauge | `namespace`, `metric_name`, `metric_type`, `hpa_owner`, `scale_target_ref_kind`, `scale_target_ref_name` |

### Grafana Daskboard

*   Get [Grafana Json](https://grafana.com/grafana/dashboards/19451-pods-metrics-server-monitor-prometheus/)

### Refer to:

*   [Pods - Metrics Server Monitor Prometheus](https://grafana.com/grafana/dashboards/8760-pods-metrics-server-monitor-prometheus/)
*   [metrics-server-monitor-prometheus](https://github.com/AdrianBalcan/metrics-server-monitor-prometheus)

### Default to enable Auto Discovery Metrics

We are configured `prometheus.io/scrape: "true"` in Service template.  
Maybe your Promethues will be automatic to get metrics and you don't need to configure anything.


#### Create Helm Package

```plaintext
helm package ./helm-chart/metrics-server-exporter --destination ./helm-chart/
helm repo index . --url  https://mrnim94.github.io/metrics-server-exporter
```