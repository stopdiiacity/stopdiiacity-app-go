include Makefile.ansible

template-generate:
	qtc -dir=templates -skipLineComments
	git add .

generate-html: template-generate
	go run ./cmd/generate-html/main.go
	git add .

swag-install:
	go get -u github.com/swaggo/swag/cmd/swag

stopdiiacity-generate-docs:
	go install github.com/swaggo/swag/cmd/swag@latest

ssh:
	ssh -t root@70.34.248.2 "cd /var/go/stopdiiacity/; bash --login"

env-up:
	mkdir -p ./.docker/volumes/go/tls-certificates
	docker-compose -f docker-compose.yml --env-file .env up -d

logs:
	docker logs stopdiiacity_go_app

app:
	docker exec -it stopdiiacity_go_app sh

env-down:
	docker-compose -f docker-compose.yml --env-file .env down

env-down-with-clear:
	docker-compose -f docker-compose.yml --env-file .env down --remove-orphans -v # --rmi=all

app-build:
	docker exec stopdiiacity_go_app go build -o /bin/stopdiiacity-server ./main.go

app-start:
	docker exec stopdiiacity_go_app stopdiiacity-server

app-stop:
	docker exec stopdiiacity_go_app pkill stopdiiacity-server || echo "stopdiiacity-server already stopped"

app-restart: app-build app-stop app-start
