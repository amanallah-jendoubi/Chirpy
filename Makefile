include .env

run:
	docker logs go-api-container -f

goose_up:
	cd /home/aj/projects/GoLang/Textio/sql/schema &&  goose postgres postgresql://postgres:$(POSTGRES_PASSWORD)@localhost:5432/postgres?sslmode=disable up


goose_down:
	cd /home/aj/projects/GoLang/Textio/sql/schema &&  goose postgres postgresql://postgres:$(POSTGRES_PASSWORD)@localhost:5432/postgres?sslmode=disable down

psql:
	docker exec -it textio-db-1 psql postgresql://postgres:$(POSTGRES_PASSWORD)@localhost:5432/postgres?sslmode=disable
sqlc:
	sqlc generate


