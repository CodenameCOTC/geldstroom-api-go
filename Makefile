dev:
	CompileDaemon -build="go build ./cmd/server/main.go" -command="./main"

migration-up:
	migrate -database mysql://root@/geldstroom  -path ./db/migrations/ up

migration-down:
	migrate -database mysql://root@/geldstroom  -path ./db/migrations/ down