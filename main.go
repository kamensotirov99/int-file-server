package main

import (
	"int-file-server/app"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.SetFormatter(&logrus.JSONFormatter{})

	go app.InitializeGrpcApp(logger)
	a := app.InitializeHttpApp(logger)
	err := a.Run()
	if err != nil {
		panic(err.Error())
	}
}
