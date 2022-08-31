lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

test:
	go test ./...
