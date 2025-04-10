package tests

import (
	"testing"

	"controller/internal/config"
	"controller/internal/controller"

	"github.com/stretchr/testify/require"
)

func TestStartAndStopService(t *testing.T) {

	err := controller.StartService(
		"test-service",
		"nginx:alpine",
		"bridge",
		map[string]string{
			"TEST_ENV": "test",
		},
	)
	require.NoError(t, err)

	err = controller.StopService("test-service")
	require.NoError(t, err)
}

func TestStartAndStopConnectorPreprocessor(t *testing.T) {
	rabbitMQURL := config.GetRabbitMQURL()
	databaseURL := config.GetDatabaseURL()

	connector := config.Connector{
		Name:        "test-connector",
		Image:       "nginx:alpine",
		Exchange:    "test-exchange",
		Queue:       "test-queue",
		RabbitMQURL: rabbitMQURL,
	}

	preprocessor := config.Preprocessor{
		Name:        "test-preprocessor",
		Image:       "nginx:alpine",
		Exchange:    "test-exchange",
		Queue:       "test-queue",
		RabbitMQURL: rabbitMQURL,
		DatabaseURL: databaseURL,
	}

	err := controller.StartConnectorAndPreprocessor(connector, preprocessor, "bridge")
	require.NoError(t, err)

	err = controller.StopConnectorAndPreprocessor(connector, preprocessor)
	require.NoError(t, err)
}

func TestUpdateServices(t *testing.T) {
	cfg := config.Config{
		Network: "bridge",
		Connectors: []config.Connector{
			{
				Name:        "test-connector-1",
				Image:       "nginx:alpine",
				Exchange:    "test-exchange",
				Queue:       "test-queue-1",
				RabbitMQURL: "amqp://guest:guest@rabbitmq:5672/",
			},
			{
				Name:        "test-connector-2",
				Image:       "nginx:alpine",
				Exchange:    "test-exchange",
				Queue:       "test-queue-2",
				RabbitMQURL: "amqp://guest:guest@rabbitmq:5672/",
			},
		},
		Preprocessors: []config.Preprocessor{
			{
				Name:        "test-preprocessor-1",
				Image:       "nginx:alpine",
				Exchange:    "test-exchange",
				Queue:       "test-queue-1",
				RabbitMQURL: "amqp://guest:guest@rabbitmq:5672/",
				DatabaseURL: "postgres://postgres:postgres@localhost:5432/heist?sslmode=disable",
			},
			{
				Name:        "test-preprocessor-2",
				Image:       "nginx:alpine",
				Exchange:    "test-exchange",
				Queue:       "test-queue-2",
				RabbitMQURL: "amqp://guest:guest@rabbitmq:5672/",
				DatabaseURL: "postgres://postgres:postgres@localhost:5432/heist?sslmode=disable",
			},
		},
	}

	controller.UpdateServices(cfg)

	for _, c := range cfg.Connectors {
		err := controller.StopService(c.Name)
		require.NoError(t, err)
	}
	for _, p := range cfg.Preprocessors {
		err := controller.StopService(p.Name)
		require.NoError(t, err)
	}
}
