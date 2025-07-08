postgres:
	docker run --name project_expiry -p 5432:5432 -e POSTGRES_USER=ropogos -e POSTGRES_PASSWORD=S3cr3t -d postgres
createdb:
	docker exec -it project_expiry createdb --username=ropogos --owner=ropogos project_expiry_db
dropdb:
	docker exec -it project_expiry dropdb -U ropogos project_expiry_db
migrateup:
	migrate -path ./migration -database "postgresql://ropogos:S3cr3t@localhost:5432/project_expiry_db?sslmode=disable" -verbose up
migratedown:
	migrate -path ./migration -database "postgresql://ropogos:S3cr3t@localhost:5432/project_expiry_db?sslmode=disable" -verbose down

.PHONY: createdb dropdb postgres migrateup migratedown