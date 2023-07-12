.PHONY: *

all: clean install lint test

install:
	go get -u ./...

lint:
	gofmt -l -e -s .
	goimports -l .

test:
	@mkdir -p var
	go test ./... -coverprofile var/cover.out

fix:
	go get -u ./...
	go mod tidy
	gofmt -s -w .
	goimports -w .

docs:
	go install golang.org/x/tools/cmd/godoc@latest
	@echo "listening on http://127.0.0.1:6060/pkg/github.com/aaronellington/zendesk-go/zendesk"
	godoc -http=127.0.0.1:6060

clean:
	git clean -Xdff
