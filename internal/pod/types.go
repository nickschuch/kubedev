package pod

import corev1 "k8s.io/api/core/v1"

const (
	ContainerName  = "app"
	ContainerLabel = "app"
)

type Params struct {
	Annotations    map[string]string `yaml:"annotations"` 
	Image          string            `yaml:"image"`
	ServiceAccount string            `yaml:"serviceAccount"`
	Env            []corev1.EnvVar   `yaml:"env"`
	Command        []string          `yaml:"command"`
}

type Mount struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}
