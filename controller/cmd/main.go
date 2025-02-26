package main

import "controller/internal/app"

const configPath = "configs/config.yaml"

func main() {
	app.Run(configPath)
}
