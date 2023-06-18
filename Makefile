.PHONY: start
start:
	docker-compose up -d --build

.PHONY: stop
stop:
	docker-compose rm -v --force --stop

.PHONY: clean
clean:
	rm -rf pgdata