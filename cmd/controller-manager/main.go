package main

/*	stopper := make(chan struct{})
	defer close(stopper)
	var kubeconfig string
	// 指定kubeconfig文件
	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")

	flag.Parse()
	controller.Run(kubeconfig, stopper)
*/

import (
	"os"

	"linkedcare.io/linkedcare/cmd/controller-manager/app"
)

func main() {
	command := app.NewControllerManagerCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
