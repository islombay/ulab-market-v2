swag:
	swag init -g api/api.go -o api/docs

run-server:
	go run cmd/main.go

swag-install:
	[ -d "api/docs" ] && echo "Directory exists" || (echo "Directory does not exist, creating"; mkdir "api/docs")
	go get -u github.com/swaggo/swag/cmd/swag@v1.16.3
	go install github.com/swaggo/swag/cmd/swag@v1.16.3

binary-start:
	./app

build: swag-install swag
	go mod tidy
	go build -o app cmd/main.go

db:
	psql -U postgres -W -h localhost -p 5432 -d ulab-market-v2

run: swag run-server
run-prod: binary-start
install: swag-install swag build