
-include env
export

api: setup
	go run cmd/api/main.go

lint: setup
	go vet ./...

setup: env
env:
	cp env.example env

