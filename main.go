package main

import (
	"GO_APP/internal/api"
	"GO_APP/config"
)

func main() {
	config := config.GetConfig()

	app := &api.App{}
	app.Initialize(config)
	app.Run(":8004")
}