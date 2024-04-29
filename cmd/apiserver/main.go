package main

import (
	"github.com/SerFiLiuZ/MEDODS/internal/app/apiserver"
	"github.com/SerFiLiuZ/MEDODS/internal/app/utils"
)

func main() {
	logger := utils.NewLogger()
	//logger.EnableDebug()

	err := apiserver.LoadEnv(logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	config := apiserver.GetConfig()

	logger.Debugf("config: %v", config)

	logger.Infof("Starting server...")

	if err := apiserver.Start(config, logger); err != nil {
		logger.Fatal(err.Error())
	}
}
