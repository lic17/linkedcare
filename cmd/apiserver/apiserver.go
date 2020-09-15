package main

import (
	"log"

	"linkedcare.io/linkedcare/cmd/apiserver/app"
)

func main() {

	cmd := app.NewAPIServerCommand()

	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
