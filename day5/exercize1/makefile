# Подключаем переменные из .env
include .env
export
# Имя Docker-сети, чтобы контейнеры видели друг друга
NETWORK=my-net

# --- Работа с миграциями (на хосте) ---
migrations-create:
	migrate create -ext sql -dir migrations -seq task_table

migrations-up:
	migrate -path migrations -database "$(DB)" up 1

migrations-down:
	migrate -path migrations -database "$(DB)" down 1

migrations-version:
	migrate -path migrations -database "$(DB)" version

# Принудительно ставим версию 1 (если база "загрязнилась")
migrations-force1:
	migrate -path migrations -database "$(DB)" force 1

# --- Локальный запуск (на хосте, для разработки) ---
go-test:
	cd features/cmd && go run main.go

# --- Docker: База данных ---
docker-create-sql:
	docker run --name postgres-run --network $(NETWORK) -e POSTGRES_PASSWORD=qwnm23 -d -v "$(PWD)/out/postgres:/var/lib/postgresql" -p 5432:5432 postgres:18-bookworm

docker-kill:
	docker kill postgres-run

docker-drop:
	docker rm postgres-run

docker-refresh:
	docker restart postgres-run

docker-stop:
	docker stop postgres-run

docker-start:
	docker start postgres-run

# --- Docker: Приложение ---
docker-create-GO:
	docker build -f dockerfile.start -t go-test .

# Запуск приложения в той же сети, что и база, с передачей .env
# не работает не пробовать
docker-run-SERV:
	docker run --rm --network my-net --env-file .env -e DB="$$DB_DOCKER" go-test 

create-net:
	docker network create my-net