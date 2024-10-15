start-rest-server:
	go run ./cmd/main.go rest

docker-build:
	docker build -t xybor/todennus-backend -f ./build/package/Dockerfile .

docker-compose-up:
	docker compose --env-file .env -f ./build/package/quick-start.yaml up -d

docker-compose-down:
	docker compose -f ./build/package/quick-start.yaml down
