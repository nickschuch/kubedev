package pod

import (
	"bufio"
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/nickschuch/kubedev/internal/log"
)

func Tail(w io.Writer, client *kubernetes.Clientset, namespace, name, container string) error {
	opts := &corev1.PodLogOptions{
		Container: container,
		Follow:    true,
	}

	body, err := client.CoreV1().Pods(namespace).GetLogs(name, opts).Stream()
	if err != nil {
		return err
	}
	defer body.Close()

	// Start reading the body and send each line to clients.
	reader := bufio.NewReader(body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			continue
		}

		log.Info(name, string(line))
	}
}
