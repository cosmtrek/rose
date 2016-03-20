BIN = $(GOPATH)/bin/rose

.PHONY: clean

clean:
	rm -f $(BIN)

build: clean
	go build

install: clean
	go install

server: install
	$(BIN)

id ?= 1
client:
	go run ./cmd/client/client.go -id=$(id)

push_message:
	go run ./cmd/pusher/push_message.go

test:
	go test -v ./protocol
