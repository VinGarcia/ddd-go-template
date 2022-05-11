
GOBIN=$(shell go env GOPATH)/bin

-include config.env
export

api: setup
	go run cmd/api/main.go

lint: setup
	@# (See staticcheck.conf to see all ignored rules)
	$(GOBIN)/staticcheck ./...
	go vet ./...

setup: config.env $(GOBIN)/staticcheck
config.env:
	cp config.env.example config.env
$(GOBIN)/staticcheck:
	GO111MODULE=off go get -u honnef.co/go/tools/cmd/staticcheck

