run:
	go run cmd/main.go

test:
	go test ./... -cover

lint:
	golangci-lint run