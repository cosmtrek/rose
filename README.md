# Rose [![Build Status](https://travis-ci.org/cosmtrek/rose.svg?branch=master)](https://travis-ci.org/cosmtrek/rose)

A super simple server supports tcp heartbeat and pushing messages to clients.

ATTENTION: Just for fun, not ready for production.

## How to play

First you have to install go develop environment.

```
cd $GO/src/github.com/cosmtrek/rose
make server
```

Open another terminal and run `make client`.

Then open another terminal and run `make push_message`.

See these three terminals output.
