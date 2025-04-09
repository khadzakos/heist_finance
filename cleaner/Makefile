.PHONY: test test-unit test-integration

test: test-unit test-integration

test-unit:
	go test -v ./tests/... -run TestNewCleaner
	go test -v ./tests/... -run TestCleaner_AddRule
	go test -v ./tests/... -run TestCleaner_RemoveRule

test-integration:
	go test -v ./tests/... -run TestCleaner_CleanData 