DB_URL=postgresql://root:123456@localhost:5432/simple_bank02?sslmode=disable

postgres:
	docker run --name postgres12-02 --network bank-network02 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:12-alpine
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
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store
.PHONY: postgres create_db drop_db migrate migrate_up migrate_up1 migrate_down migrate_down1 sqlc test server mock