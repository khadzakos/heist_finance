.PHONY: test test-unit test-integration

test: test-unit test-integration

test-unit:
	go test -v ./tests/... -run TestNewConnector
	go test -v ./tests/... -run TestConnector_Connect
	go test -v ./tests/... -run TestConnector_Disconnect

test-integration:
	go test -v ./tests/... -run TestConnector_SendData 