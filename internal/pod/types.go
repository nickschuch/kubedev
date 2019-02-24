package pod

import corev1 "k8s.io/api/core/v1"

const ContainerName = "app"

type Params struct {
	Name           string          `yaml:"name"`
	Image          string          `yaml:"image"`
	ServiceAccount string          `yaml:"serviceAccount"`
	Env            []corev1.EnvVar `yaml:"env"`
	Command        []string        `yaml:"command"`
}

type Mount struct {
	Name   string `yaml:"name"`
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}
