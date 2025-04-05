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
test:
	@for module in auth common gateway payments; do \
		echo "Testing $$module module..."; \
		go test -v ./$$module/... || exit 1; \
	done

.PHONY: help
help:
	@echo "Available commands:"
	@echo " make proto           - Generate protobuf files for all services"
	@echo " make run-auth        - Run the auth service"
	@echo " make run-payments    - Run the payments service"
	@echo " make run-all         - Run all services (in separate terminals)"


SERVICES := auth payments
.PHONY: proto
proto: $(SERVICES)

.PHONY: $(SERVICES)
$(SERVICES):
	rm -f $@/pb/*.go || true
	protoc --proto_path=$@/proto \
	--go_out=$@/pb --go_opt=paths=source_relative \
	--go-grpc_out=$@/pb --go-grpc_opt=paths=source_relative \
	$@/proto/*.proto

.PHONY: run-auth
run-auth:
	@echo "Starting auth service..."
	cd auth && go run main.go

.PHONY: run-payments
run-payments:
	@echo "Starting payments service..."
	cd payments && go run main.go

.PHONY: run-all
run-all:
	@echo "Starting all services..."
	@for service in $(SERVICES); do \
		if [ "$(shell uname)" = "Linux" ]; then \
			gnome-terminal -- make run-$$service; \
		elif [ "$(shell uname)" = "Darwin" ]; then \
			osascript -e 'tell app "Terminal" to do script "cd $(PWD) && make run-'$$service'"'; \
		else \
			echo "Unsupported OS for terminal automation"; \
			exit 1; \
		fi; \
	done
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test