include .env

run:
	docker logs go-api-container -f

goose_up:
	cd /home/aj/projects/GoLang/Matrix_Chat/sql/schema &&  goose postgres postgresql://postgres:$(POSTGRES_PASSWORD)@localhost:5432/postgres?sslmode=disable up


goose_down:
	cd /home/aj/projects/GoLang/Matrix_Chat/sql/schema &&  goose postgres postgresql://postgres:$(POSTGRES_PASSWORD)@localhost:5432/postgres?sslmode=disable down

psql:
	docker exec -it go-db-container psql postgresql://postgres:$(POSTGRES_PASSWORD)@localhost:5432/postgres?sslmode=disable
sqlc:
	sqlc generate


