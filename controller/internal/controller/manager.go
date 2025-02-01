package controller

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ConnectorConfig struct {
	Name    string
	Image   string
	WsUrl   string
	Tickers []string
}

func StartConnector(cfg ConnectorConfig) {
	cli, _ := client.NewClientWithOpts(client.FromEnv)

	envVars := []string{
		"WS_URL=" + cfg.WsUrl,
		"TICKERS=" + strings.Join(cfg.Tickers, ","),
	}

	_, err := cli.ContainerCreate(context.Background(),
		&container.Config{
			Image: cfg.Image,
			Env:   envVars,
		}, nil, nil, nil, cfg.Name)

	if err != nil {
		log.Printf("Ошибка запуска %s: %v", cfg.Name, err)
		return
	}

	cli.ContainerStart(context.Background(), cfg.Name, container.StartOptions{})
	log.Printf("Запущен коннектор: %s (Tickers: %v)", cfg.Name, cfg.Tickers)

	// Добавляем чтение логов после запуска
	go func() {
		reader, err := cli.ContainerLogs(context.Background(), cfg.Name, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		})
		if err != nil {
			log.Printf("Ошибка при получении логов %s: %v", cfg.Name, err)
			return
		}
		defer reader.Close()

		buf := make([]byte, 1024)
		for {
			n, err := reader.Read(buf)
			if err != nil {
				break
			}
			log.Print(string(buf[:n]))
		}
	}()
}

func StopConnector(name string) {
	cli, _ := client.NewClientWithOpts(client.FromEnv)
	cli.ContainerStop(context.Background(), name, container.StopOptions{})
	cli.ContainerRemove(context.Background(), name, container.RemoveOptions{})
	log.Printf("Остановлен коннектор: %s", name)
}
