template-generate:
	qtc -dir=templates -skipLineComments
	git add .

generate-html: template-generate
	go run ./cmd/main.go
	git add .

swag-install:
	go get -u github.com/swaggo/swag/cmd/swag

stopdiiacity-generate-docs:
	swag init -o apidocs -g main.go

run:
	browse http://localhost:8080
	PORT=8080 go run main.go
