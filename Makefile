include Makefile.ansible

template-generate:
	qtc -dir=templates -skipLineComments
	git add .

generate-html: template-generate
	go run ./cmd/main.go
	git add .

swag-install:
	go get -u github.com/swaggo/swag/cmd/swag

stopdiiacity-generate-docs:
	go install github.com/swaggo/swag/cmd/swag@latest

ssh:
	ssh -t root@70.34.251.121 "cd /var/go/stopdiiacity/; bash --login"

run:
	mkdir -p ./.docker/volumes/go/tls-certificates
	browse http://localhost
	PORT="80" \
 		TLS_CERTIFICATES_DIR="./.docker/volumes/go/tls-certificates" \
 		HOSTS="stopdiiacity.u8hub.com" \
 		go run main.go
