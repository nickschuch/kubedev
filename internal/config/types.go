package config

import "github.com/nickschuch/kubedev/internal/pod"

type File struct {
	Namespace string                `yaml:"namespace"`
	Mounts    map[string]pod.Mount  `yaml:"mounts"`
	Pods      map[string]pod.Params `yaml:"pods"`
}
