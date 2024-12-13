.PHONY: stop kill startCompose startApiService startAuthService startSmtpService

all: startCompose startApiService startAuthService startSmtpService

stop:
	docker stop $(shell docker ps -aq)

kill:
	docker rm $(shell docker ps -aq)

startCompose:
	docker-compose -f ./docker/compose/docker-compose.yml up -d

startApiService:
	go build -o ./api-service/cmd/main ./api-service/cmd/main.go
	pkill -f './api-service/cmd/main'
	./api-service/cmd/main &

startAuthService:
	go build -o ./auth-service/cmd/main ./auth-service/cmd/main.go 
	pkill -f './auth-service/cmd/main'
	./auth-service/cmd/main &

startSmtpService:
	go build -o ./smtp-service/cmd/main ./smtp-service/cmd/main.go 
	pkill -f './smtp-service/cmd/main'
	./smtp-service/cmd/main &

