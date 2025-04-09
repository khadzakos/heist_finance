.PHONY: test test-unit test-integration

test: test-unit test-integration

test-unit:
	go test -v ./tests/... -run TestNewProcessor
	go test -v ./tests/... -run TestProcessor_AddTransform
	go test -v ./tests/... -run TestProcessor_RemoveTransform

test-integration:
	go test -v ./tests/... -run TestProcessor_ProcessData 