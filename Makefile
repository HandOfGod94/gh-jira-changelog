APP_NAME=gh-jira-changelog

all: clean fmt tidy build install

clean:
	rm -f $(APP_NAME)

install:
	gh extension install .

quality-check:
	staticcheck ./...
	gocyclo -over 15 -ignore "testdata/" .
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
