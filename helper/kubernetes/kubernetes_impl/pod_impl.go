package kubernetes_impl

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"metrics-server-exporter/log"
)

func (kc *KubeConfiguration) DetectPodOs(namespace, podName string) (string, error) {
	config, err := kc.accessKubernetes()
	if err != nil {
		log.Error(err.Error())
	}

	// Create a clientset to interact with the Kubernetes API
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err.Error())
	}

	// Get a list of all Pods in all namespaces
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		return "", err
	}

	if os, exists := pod.Spec.NodeSelector["kubernetes.io/os"]; exists {
		return os, nil // returns "linux" or "windows"
	}

	return "not_specified", nil // OS nodeSelector is not configured
}

// Function to get the owner (Deployment, StatefulSet, DaemonSet) based on Pod name and namespace
func (kc *KubeConfiguration) GetPodOwner(namespace, podName string) (ownerKind, ownerName string, err error) { //NOSONAR - Refactor this method to reduce its Cognitive Complexity
	// Load kubeconfig
	config, err := kc.accessKubernetes()
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	// Get Pod
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	// Check owner references of the Pod
	for _, ownerReference := range pod.OwnerReferences {
		switch ownerReference.Kind {
		case "ReplicaSet":
			// Get the ReplicaSet
			rs, err := clientset.AppsV1().ReplicaSets(namespace).Get(context.TODO(), ownerReference.Name, metav1.GetOptions{})
			if err != nil {
				return "", "", err
			}
			// Check owner references of the ReplicaSet to find the Deployment
			for _, rsOwnerReference := range rs.OwnerReferences {
				if rsOwnerReference.Kind == "Deployment" {
					return "Deployment", rsOwnerReference.Name, nil // Found the Deployment
				}
			}
		case "StatefulSet":
			return "StatefulSet", ownerReference.Name, nil
		case "DaemonSet":
			return "DaemonSet", ownerReference.Name, nil
		}
	}
	log.Debug(fmt.Errorf("no Deployment, StatefulSet, or DaemonSet found for Pod %s in namespace %s", podName, namespace))
	return "", "", nil
}
