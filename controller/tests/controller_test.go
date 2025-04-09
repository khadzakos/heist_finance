package tests

import (
	"testing"

	"controller/internal/controller"

	"github.com/stretchr/testify/require"
)

func TestStartAndStopService(t *testing.T) {
	err := controller.StartService(
		"test-service",
		"hello-world",
		"bridge",
		map[string]string{
			"TEST_ENV": "test",
		},
	)
	require.NoError(t, err)

	err = controller.StopService("test-service")
	require.NoError(t, err)
}
