.PHONY: start
start:
	docker-compose up -d --build

.PHONY: stop
stop:
	docker-compose down

.PHONY: clean
clean:
	rm -rf pgdata

.PHONY: lint
lint:
	docker run \
		-t \
		--rm \
		-v $(pwd):/app \
		-w /app golangci/golangci-lint:v1.53.3 \
		golangci-lint run -v