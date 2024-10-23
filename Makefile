start-rest:
	go run ./cmd/main.go rest

start-swagger:
	go run ./cmd/main.go swagger

docker-build:
	docker build -t xybor/todennus-backend -f ./build/package/Dockerfile .

swagger-gen:
	swag init --dir ./adapter/rest/ -g app.go
