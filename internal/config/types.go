package config

import "github.com/nickschuch/kubedev/internal/pod"

type File struct {
	Namespace string            `yaml:"namespace"`
	Labels    map[string]string `yaml:"labels"`
	Mounts    []pod.Mount       `yaml:"mounts"`
	Pods      []pod.Params      `yaml:"pods"`
}
