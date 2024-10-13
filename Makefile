gen_postgres_migration:
	./build/gen_migrations.sh postgres $(name)

down_postgres_migration:
	./build/down_migration.sh postgres 1

start-rest-server:
	go run ./cmd/main.go rest

docker-build:
	docker build -t xybor/todennus-backend -f ./build/package/Dockerfile .

docker-compose-up:
	docker compose --env-file .env -f ./build/package/quick-start.yaml up -d

docker-compose-down:
	docker compose -f ./build/package/quick-start.yaml down
