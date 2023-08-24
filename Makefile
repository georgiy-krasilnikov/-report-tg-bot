# SERVICE_NAME=report-bot
# PORT=8080
# PWD=$(dir $(abspath $(lastword $(MAKEFILE_LIST))))


# up:
# 	docker build -t $(SERVICE_NAME) .
# 	docker run -v ./docs:$(PWD)docs -p $(PORT):$(PORT) --env-file=.env  $(SERVICE_NAME)
up:
	go run main.go