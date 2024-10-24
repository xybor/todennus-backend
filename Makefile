start-rest:
	go run ./cmd/main.go rest

start-grpc:
	go run ./cmd/main.go grpc

start-swagger:
	go run ./cmd/main.go swagger

docker-build:
	docker build -t xybor/todennus-backend -f ./build/package/Dockerfile .

swagger-gen:
	swag init --dir ./adapter/rest/ -g app.go

proto-gen:
	rm -rf ./adapter/grpc/gen/* && \
	protoc --go_out=./adapter/grpc/gen --go_opt=paths=source_relative \
    	--go-grpc_out=./adapter/grpc/gen --go-grpc_opt=paths=source_relative \
    	-I=../todennus-proto/ -I=../todennus-proto/dto/  -I=../todennus-proto/dto/resource/ \
		../todennus-proto/dto/resource/*.proto ../todennus-proto/dto/*.proto ../todennus-proto/*.proto
