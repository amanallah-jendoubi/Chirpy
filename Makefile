include .env

run:
	docker logs go-api-container -f

goose_up:
	cd sql/schema &&  goose postgres postgresql://postgres:$(POSTGRES_PASSWORD)@localhost:5432/postgres?sslmode=disable up


goose_down:
	cd sql/schema &&  goose postgres postgresql://postgres:$(POSTGRES_PASSWORD)@localhost:5432/postgres?sslmode=disable down

psql:
	docker exec -it go-db-container psql postgresql://postgres:$(POSTGRES_PASSWORD)@localhost:5432/postgres?sslmode=disable
sqlc:
	sqlc generate


