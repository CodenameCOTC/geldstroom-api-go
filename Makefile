dev:
	CompileDaemon -build="go build ./cmd/server/main.go" -command="./main"

migration-up:
	migrate -database mysql://root@/geldstroom  -path ./db/migrations/ up

migration-down:
	migrate -database mysql://root@/geldstroom  -path ./db/migrations/ down

backup-db:
	docker exec db /usr/bin/mysqldump -u root geldstroom > backup.sql	

restore-db:
	cat backup.sql | docker exec -i db /usr/bin/mysql -u root geldstroom	