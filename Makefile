lint:
	golangci-lint run

fmt:
	gofmt -s -w .

run:
	go run main.go