swag:
	swag init -g api/api.go -o api/docs

run-server:
	go run cmd/main.go

run: swag run-server