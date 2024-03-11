package model

type HPAUtilizations []HPAUtilization

type HPAUtilization struct {
	MetricName         string
	MetricType         string
	HPAOwner           string
	ScaleTargetRefKind string
	ScaleTargetRefName string
	AverageUtilization float64
}
