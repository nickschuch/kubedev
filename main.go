package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/nickschuch/kubedev/internal/config"
	"github.com/nickschuch/kubedev/internal/log"
	"github.com/nickschuch/kubedev/internal/pod"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	cmdFile   = kingpin.Flag("file", "File to load").Default("kubedev.yml").String()
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
	for _, params := range file.Pods {
		go func(params pod.Params) {
			log.Infoln(params.Name, "Starting")

			err := pod.Run(kubeclient, file.Namespace, file.Labels, params, file.Mounts)
			if err != nil {
				log.Error(params.Name, "Pod run failed:", err)
			}

			log.Infoln(params.Name, "Waiting to become running")

			err = pod.Wait(kubeclient, file.Namespace, params.Name)
			if err != nil {
				log.Error(params.Name, "Pod wait failed:", err)
			}

			log.Infoln(params.Name, "Starting log stream")

			err = pod.Tail(os.Stdout, kubeclient, file.Namespace, params.Name, pod.ContainerName)
			if err != nil {
				log.Error(params.Name, "Pod log stream failed:", err)
			}
		}(params)
	}

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})

	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan

		fmt.Println("Received an interrupt, terminating Pods...")

		for _, pod := range file.Pods {
			log.Infoln(pod.Name, "Terminating")

			err := kubeclient.CoreV1().Pods(file.Namespace).Delete(pod.Name, &metav1.DeleteOptions{})
			if err != nil {
				log.Error(pod.Name, "Failed to terminate Pod:", err)
			}
		}

		close(cleanupDone)
	}()

	<-cleanupDone
}
