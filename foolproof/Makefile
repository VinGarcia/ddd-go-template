
-include config.env
export

api: setup
	go run cmd/api/main.go

lint: setup
	go vet ./...

setup: config.env
config.env:
	cp config.env.example config.env

