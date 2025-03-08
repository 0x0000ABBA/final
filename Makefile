test:
	go test ./...
docker-run:
	docker compose up -d --build
build:
	go build ./cmd/main/main.go
run:
	go run ./cmd/main/main.go
lint:
	golangci-lint run