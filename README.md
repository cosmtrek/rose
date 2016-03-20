# Rose [![Build Status](https://travis-ci.org/cosmtrek/rose.svg?branch=master)](https://travis-ci.org/cosmtrek/rose)

A super simple server supports tcp heartbeat and pushing messages to clients.

ATTENTION: Just for fun, not ready for production.

## How to play

First you need to install Go 1.6.

```
go get github.com/cosmtrek/rose
cd $GOPATH/src/github.com/cosmtrek/rose
make server
```

Open another terminal and run `make client`. You can run `make client id=2` to specify which client to be connecting the server and so on.

In addition to that, open another terminal and run `make push_message`.

Check these terminals output.
