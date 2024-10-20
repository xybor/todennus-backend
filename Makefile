start-rest:
	go run ./cmd/main.go rest

docker-build:
	docker build -t xybor/todennus-backend -f ./build/package/Dockerfile .
