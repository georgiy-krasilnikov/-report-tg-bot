SERVICE_NAME=report-bot
PORT=8080

up: 
	docker build -t $(SERVICE_NAME) .
	docker run -p $(PORT):$(PORT)  --env-file=.env  $(SERVICE_NAME)