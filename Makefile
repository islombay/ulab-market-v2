swag:
	swag init -g api/api.go -o api/docs

run-server:
	go run cmd/main.go

swag-install:
	[ -d "api/docs" ] || mkdir api/docs
	go get -u github.com/swaggo/swag/cmd/swag
	go install github.com/swaggo/swag/cmd/swag@v1.16.3

binary-start:
	./app

build:
	go build -o app cmd/main.go

run: swag run-server
run-prod: swag binary-start
install: swag-install build