package model

import "time"

type PodUsage []struct {
	Metadata struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Labels            struct {
			App             string `json:"app"`
			PodTemplateHash string `json:"pod-template-hash"`
			Release         string `json:"release"`
		} `json:"labels"`
	} `json:"metadata"`
	Timestamp  time.Time `json:"timestamp"`
	Window     string    `json:"window"`
	Containers []struct {
		Name  string `json:"name"`
		Usage struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"usage"`
	} `json:"containers"`
}
