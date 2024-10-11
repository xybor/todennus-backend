gen_postgres_migration:
	./build/gen_migrations.sh postgres $(name)

down_postgres_migration:
	./build/down_migration.sh postgres 1

start-rest-server:
	go run ./cmd/main.go rest
