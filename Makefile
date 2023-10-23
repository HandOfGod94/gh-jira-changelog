APP_NAME=gh-jira-changelog
EXTENSION_NAME=jira-changelog

all: clean fmt tidy build install quality-check

clean:
	rm -f $(APP_NAME)
	gh extension remove $(EXTENSION_NAME) &2>/dev/null; true

install:
	go install
	go build -o $(APP_NAME) && gh extension install .

quality-check:
	staticcheck ./...
	gocyclo -over 15 .
	gocognit -over 15 .

test:
	go test -timeout 60s -cover -count=1 -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

tidy:
	go mod tidy -v

run: dev
	$(APP_NAME)

dev:
	go build -v -o ouAPP_NAME) main.go

build:
	go build -v -ldflags "-w" -o $(APP_NAME) main.go

fmt:
	go fmt ./...
	gofmt -s -w .

setup-tools:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/uudashr/gocognit/cmd/gocognit@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/mcubik/goverreport@latest
	go install github.com/vektra/mockery/v2@v2.32.0

sync-docs:
	mdsh
