run:
	go run cmd/main.go
db:
	docker-compose up
image:
	docker build -t apod .

.PHONY: run