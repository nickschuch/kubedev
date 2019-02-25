package pod

import (
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/pkg/errors"
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

	_, err = io.Copy(w, body)
	if err != nil {
		return errors.Wrap(err, "failed to copy logs")
	}

	return nil
}
