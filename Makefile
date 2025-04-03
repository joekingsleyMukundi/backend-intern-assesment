postgres:
	docker run --name Pdpostgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mypassword -d postgres:17-alpine
createdb:
	docker exec -it  Pdpostgres createdb --username=root --owner=root payd
dropdb:
	docker exec -it Pdpostgres dropdb payd
migrateup:
	migrate -path db/migrations -database "postgresql://root:mypassword@localhost:5432/payd?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://root:mypassword@localhost:5432/payd?sslmode=disable" -verbose down
sqlc:
	cd common && sqlc generate
	.PHONY: test
test:
	@for module in auth common gateway payments; do \
		echo "Testing $$module module..."; \
		go test -v ./$$module/... || exit 1; \
	done
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test