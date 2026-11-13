include .env

postgres-run:
	docker run --name postgres-run -e POSTGRES_PASSWORD=qwnm23 -d -v "/$PWD/out/postgres:/var/lib/postgresql" -p 5432:5432 postgres:18-bookworm

porstgres-kill:
	docker kill postgres-run 

docker-kill-post:
	docker rm -f postgres-run
go-run:
	go run ./cmd/server/

migration-create:
	migrate create -ext sql -dir migrations -seq create_users_table


migration-Up:
	migrate -path migrations -database "$(DB)" up 1

migration-Down:
	migrate -path migrations -database "$(DB)" down 1