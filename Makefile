swag:
	swag init -g api/api.go -o api/docs

run-server:
	go run cmd/main.go

swag-install:
	go get -u github.com/swaggo/swag/cmd/swag
	go install github.com/swaggo/swag/cmd/swag@v1.16.3
run: swag-install swag run-server