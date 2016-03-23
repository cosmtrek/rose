BIN := $(GOPATH)/bin/rose

.PHONY: clean

clean:
	rm -f $(BIN)

build: clean
	go build

install: clean
	go install

dev: install
	sh script.sh
	$(BIN)

production: install
	sh script.sh
	$(BIN) > rose.production.log &

id ?= 1
client:
	go run ./cmd/client/client.go -id=$(id)

push_message:
	go run ./cmd/pusher/push_message.go

test:
	go test -v ./protocol
