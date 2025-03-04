test:
	go test ./...
docker-build:
	docker compose build --no-cache
docker-run:
	docker compose up -d
build:
	go build ./cmd/main/main.go
run:
	go run ./cmd/main/main.go
lint:
	golangci-lint run