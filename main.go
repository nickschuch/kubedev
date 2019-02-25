package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/nickschuch/kubedev/internal/config"
	"github.com/nickschuch/kubedev/internal/log"
	"github.com/nickschuch/kubedev/internal/pod"
)

var (
	cmdFile   = kingpin.Flag("file", "File to load").Default("kubedev.yml").String()
	cmdFilter = kingpin.Flag("filter", "Filter by program name eg. app1,app2").String()
	cmdConfig = kingpin.Arg("kubeconfig", "Path to the Kubeconfig file").Envar("KUBECONFIG").String()
)

func main() {
	kingpin.Parse()

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", *cmdConfig)
	if err != nil {
		panic(err)
	}

	kubeclient, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		panic(err)
	}

	file, err := config.LoadFromFile(*cmdFile)
	if err != nil {
		panic(err)
	}

	// This workgroup is for adding Pods to our execution.
	for name, params := range file.Pods {
		if filter(*cmdFilter, name) {
			continue
		}

		go func(name string, params pod.Params) {
			log.Infoln(name, "Starting")

			err := pod.Run(kubeclient, file.Namespace, name, params, file.Mounts)
			if err != nil {
				log.Error(name, "Pod run failed:", err)
			}

			log.Infoln(name, "Waiting to become running")

			err = pod.Wait(kubeclient, file.Namespace, name)
			if err != nil {
				log.Error(name, "Pod wait failed:", err)
			}

			log.Infoln(name, "Starting log stream")

			err = pod.Tail(os.Stdout, kubeclient, file.Namespace, name, pod.ContainerName)
			if err != nil {
				log.Error(name, "Pod log stream failed:", err)
			}
		}(name, params)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Received an interrupt, terminating Pods...")

	for name := range file.Pods {
		if filter(*cmdFilter, name) {
			continue
		}

		log.Infoln(name, "Terminating")

		err := kubeclient.CoreV1().Pods(file.Namespace).Delete(name, &metav1.DeleteOptions{})
		if err != nil {
			log.Error(name, "Failed to terminate Pod:", err)
		}
	}
}

func filter(list, name string) bool {
	if list == "" {
		return false
	}

	// Can we find the application in the list.
	for _, a := range strings.Split(list, ",") {
		if a == name {
			return false
		}
	}

	return true
}
