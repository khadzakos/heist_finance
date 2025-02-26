package main

import "controller/internal/app"

const configPath = "/app/configs/config.yaml"

func main() {
	app.Run(configPath)
}
