.PHONY: test test-unit test-integration

test: test-unit test-integration

test-unit:
	go test -v ./tests/... -run TestNewController
	go test -v ./tests/... -run TestController_AddHandler
	go test -v ./tests/... -run TestController_RemoveHandler

test-integration:
	go test -v ./tests/... -run TestController_ProcessRequest 