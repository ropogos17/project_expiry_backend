postgres:
	docker run --name project_expiry -p 5432:5432 --env-file app.env -d postgres
createdb:
	docker exec -it project_expiry createdb --username=ropogos --owner=ropogos project_expiry_db
dropdb:
	docker exec -it project_expiry dropdb -U ropogos project_expiry_db
migrateup:
	go run ./cmd/migrate up
migratedown:
	go run ./cmd/migrate down
sqlc:
	sqlc generate
test:
	go test -v -cover -count=1 ./...
server:
	go run main.go

.PHONY: createdb dropdb postgres migrateup migratedown sqlc test server