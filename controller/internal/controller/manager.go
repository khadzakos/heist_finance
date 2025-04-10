package controller

import (
	"context"
	"controller/internal/config"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var runningConnectors = make(map[string]bool)
var mu sync.Mutex

// Создание клиента Docker
func newDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

// Проверка существующего контейнера
func containerExists(cli *client.Client, name string) (bool, string, error) {
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return false, "", err
	}

	for _, c := range containers {
		for _, containerName := range c.Names {
			if containerName == "/"+name {
				return true, c.ID, nil
			}
		}
	}
	return false, "", nil
}

// Функция для проверки доступности RabbitMQ
func waitForRabbitMQ(host string, port string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), time.Second)
		if err == nil {
			conn.Close()
			log.Printf("RabbitMQ доступен после %d попыток\n", tries+1)
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("timeout waiting for RabbitMQ to become available")
}

// Запуск сервиса (контейнера)
func StartService(name, image, network string, env map[string]string) error {
	cli, err := newDockerClient()
	if err != nil {
		log.Printf("Ошибка создания клиента Docker: %v\n", err)
		return err
	}
	defer cli.Close()

	exists, containerID, err := containerExists(cli, name)
	if err != nil {
		log.Printf("Ошибка проверки контейнера %s: %v\n", name, err)
		return err
	}

	if exists {
		log.Printf("Контейнер %s уже существует, удаляем...", name)
		err = cli.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true})
		if err != nil {
			log.Printf("Ошибка удаления контейнера %s: %v\n", name, err)
			return err
		}
	}

	var envVars []string
	for key, value := range env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	resp, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
			Env:   envVars,
		},
		&container.HostConfig{
			NetworkMode: container.NetworkMode(network),
		},
		nil,
		nil,
		name,
	)

	if err != nil {
		log.Printf("Ошибка создания контейнера %s: %v\n", name, err)
		return err
	}

	err = cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
	if err != nil {
		log.Printf("Ошибка запуска контейнера %s: %v\n", name, err)
		return err
	}
	if name == "rabbitmq" {
		if err := waitForRabbitMQ("rabbitmq", "5672", 60*time.Second); err != nil {
			log.Printf("Ошибка ожидания RabbitMQ: %v\n", err)
			return err
		}
	}

	log.Printf("Запущен сервис: %s\n", name)
	return nil
}

// Остановка и удаление контейнера
func StopService(name string) error {
	cli, err := newDockerClient()
	if err != nil {
		log.Printf("Ошибка создания клиента Docker: %v\n", err)
		return err
	}
	defer cli.Close()

	exists, containerID, err := containerExists(cli, name)
	if err != nil {
		log.Printf("Ошибка проверки контейнера %s: %v\n", name, err)
		return err
	}

	if !exists {
		log.Printf("Контейнер %s не найден, пропускаем остановку.\n", name)
		return nil
	}

	log.Printf("Остановка контейнера %s...", name)
	if err := cli.ContainerStop(context.Background(), containerID, container.StopOptions{}); err != nil {
		log.Printf("Ошибка остановки контейнера %s: %v\n", name, err)
		return err
	}

	log.Printf("Удаление контейнера %s...", name)
	if err := cli.ContainerRemove(context.Background(), containerID, container.RemoveOptions{}); err != nil {
		log.Printf("Ошибка удаления контейнера %s: %v\n", name, err)
		return err
	}

	log.Printf("Остановлен сервис: %s\n", name)
	return nil
}

// Запуск связки connector + preprocessor
func StartConnectorAndPreprocessor(c config.Connector, p config.Preprocessor, network string) error {
	if err := waitForRabbitMQ("rabbitmq", "5672", 30*time.Second); err != nil {
		return fmt.Errorf("RabbitMQ недоступен: %v", err)
	}

	err := StartService(c.Name, c.Image, network, map[string]string{
		"EXCHANGE":     c.Exchange,
		"QUEUE":        c.Queue,
		"RABBITMQ_URL": c.RabbitMQURL,
	})
	if err != nil {
		return err
	}

	err = StartService(p.Name, p.Image, network, map[string]string{
		"QUEUE":        p.Queue,
		"EXCHANGE":     p.Exchange,
		"RABBITMQ_URL": p.RabbitMQURL,
		"DATABASE_URL": p.DatabaseURL,
	})
	if err != nil {
		return err
	}

	return nil
}

// Остановка связки connector + preprocessor
func StopConnectorAndPreprocessor(c config.Connector, p config.Preprocessor) error {
	if err := StopService(c.Name); err != nil {
		return err
	}
	if err := StopService(p.Name); err != nil {
		return err
	}
	return nil
}

// Обновление сервисов (автоматический запуск/остановка)
func UpdateServices(newConfig config.Config) {
	mu.Lock()
	defer mu.Unlock()

	current := make(map[string]bool)

	for _, c := range newConfig.Connectors {
		for _, p := range newConfig.Preprocessors {
			if c.Queue == p.Queue {
				current[c.Name] = true
				if !runningConnectors[c.Name] {
					go StartConnectorAndPreprocessor(c, p, newConfig.Network)
				}
			}
		}
	}

	// TODO: остановка сервисов, которых больше нет в конфиге - починить!
	// Остановка сервисов, которых больше нет в конфиге
	// for name := range runningConnectors {
	// 	if !current[name] {
	// 		go StopService(name)
	// 	}
	// }

	runningConnectors = current
}
