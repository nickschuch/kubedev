package pod

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Wait for a pod to become available.
func Wait(clientset *kubernetes.Clientset, namespace, name string) error {
	watcher, err := clientset.CoreV1().Pods(namespace).Watch(metav1.ListOptions{})
	if err != nil {
		return err
	}
	defer watcher.Stop()

	for {
		event := <-watcher.ResultChan()

		if event.Object == nil {
			continue
		}

		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			continue
		}

		if pod.Name != name {
			continue
		}

		if pod.Status.Phase == corev1.PodRunning {
			return nil
		}

		if pod.Status.Phase == corev1.PodFailed {
			return fmt.Errorf("Pod status: %s", corev1.PodFailed)
		}

		if pod.Status.Phase == corev1.PodUnknown {
			return fmt.Errorf("Pod status: %s", corev1.PodUnknown)
		}
	}
}
