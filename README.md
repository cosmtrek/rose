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

## Configuration

There two ways to customize some arguments, one is copying `.rose.toml` to your home directory where Rose will check when starting, the other is running `rose -server_host=localhost -server_port=3333 -socket_timeout=300`.
