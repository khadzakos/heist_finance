package controller

import (
	"github.com/docker/docker/client"
)

func MonitorConnectors() {
	cli, _ := client.NewClientWithOpts(client.FromEnv)
	// TODO: Добавить мониторинг коннекторов
}
