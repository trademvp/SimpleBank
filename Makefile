DB_URL=postgresql://root:123456@localhost:5432/simple_bank02?sslmode=disable

network:
	docker network create bank-network
postgres:
	docker run --name postgres12-02 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:12-alpine
create_db:
	docker exec -it postgres12-02 createdb --username=root --owner=root simple_bank02
drop_db:
	docker exec -it postgres12-02 dropdb simple_bank02
migrate:
	 migrate create -ext sql -dir db/migration -seq add_users
migrate_up:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migrate_up1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
migrate_down:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
migrate_down1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
db_docs:
	dbdocs build doc/db.dbml
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store
proto:
	rm -f pb/*.go
	statik -src=./doc/swagger -dest=./doc
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl
redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine
redis_ping:
	docker exec -it redis redis-cli ping

.PHONY: postgres create_db drop_db migrate migrate_up migrate_up1 migrate_down migrate_down1 sqlc test server mock proto evans db_docs db_schema redis redis_ping