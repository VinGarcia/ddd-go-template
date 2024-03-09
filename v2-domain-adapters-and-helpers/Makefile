
GOBIN=$(shell go env GOPATH)/bin
path=./...

-include config.env
export

api: setup
	go run cmd/api/main.go

test: setup
	@# First try to pull the postgres image used on the tests
	@# otherwise the tests will fail with a timeout error since
	@# this might take a long time to download.
	@docker pull postgres:14.8
	$(GOBIN)/richgo test $(path) $(args)

lint: setup
	@# (See staticcheck.conf to see all ignored rules)
	$(GOBIN)/staticcheck ./...
	go vet ./...

setup: config.env $(GOBIN)/richgo $(GOBIN)/staticcheck

config.env:
	cp config.env.example config.env

$(GOBIN)/richgo:
	go install github.com/kyoh86/richgo@latest

$(GOBIN)/staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest

