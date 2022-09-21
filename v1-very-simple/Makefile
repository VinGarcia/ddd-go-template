
GOBIN=$(shell go env GOPATH)/bin
path=./...

-include config.env
export

api: setup
	go run cmd/api/main.go

test: setup
	$(GOBIN)/richgo test $(path) $(args)

lint: setup
	@# (See staticcheck.conf to see all ignored rules)
	$(GOBIN)/staticcheck ./...
	go vet ./...

setup: config.env $(GOBIN)/richgo $(GOBIN)/staticcheck
config.env:
	cp config.env.example config.env

$(GOBIN)/richgo:
	go get github.com/kyoh86/richgo

$(GOBIN)/staticcheck:
	GO111MODULE=off go get -u honnef.co/go/tools/cmd/staticcheck

