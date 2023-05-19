package main

import (
	"GO_APP/config"
	api "GO_APP/internal/delivery"
)

func main() {
	config := config.GetConfig()

	app := &api.App{}
	app.Init(config)
	app.Run(":8004")
}
