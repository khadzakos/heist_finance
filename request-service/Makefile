.PHONY: test test-unit test-integration

test: test-unit test-integration

test-unit:
	go test -v ./tests/... -run TestNewRequestService
	go test -v ./tests/... -run TestRequestService_AddEndpoint
	go test -v ./tests/... -run TestRequestService_RemoveEndpoint

test-integration:
	go test -v ./tests/... -run TestRequestService_HandleRequest 