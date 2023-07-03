.PHONY: start
start:
	docker-compose up -d --build

.PHONY: stop
stop:
	docker-compose down

.PHONY: clean
clean:
	rm -rf pgdata

.PHONY: run
run:
	go run cmd/shortener/main.go -d "postgres://postgres:password@localhost:5432/shorten?sslmode=disable"

.PHONY: lint
lint:
	docker run \
		-t \
		--rm \
		-v $(pwd):/app \
		-w /app golangci/golangci-lint:v1.53.3 \
		golangci-lint run -v