package pod

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Run(client *kubernetes.Clientset, namespace, name string, params Params, mounts map[string]Mount) error {
	container := corev1.Container{
		Name:    ContainerName,
		Image:   params.Image,
		Env:     params.Env,
		Command: params.Command,
	}

	for mountName, mountValue := range mounts {
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			Name:      mountName,
			MountPath: mountValue.Target,
		})
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Annotations: params.Annotations,
			Labels: map[string]string{
				ContainerLabel: name,
			},
		},
		Spec: corev1.PodSpec{
			Containers:         []corev1.Container{container},
			ServiceAccountName: params.ServiceAccount,
			RestartPolicy:      corev1.RestartPolicyNever,
		},
	}

	for mountName, mountValue := range mounts {
		pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
			Name: mountName,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: mountValue.Source,
				},
			},
		})
	}

	pod, err := client.CoreV1().Pods(namespace).Create(pod)
	if err != nil {
		return errors.Wrap(err, "cannot create Pod")
	}

	return nil
}
