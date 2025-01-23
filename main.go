package main

import (
	"github.com/jeffscottbrown/applemusic/logging"
	"github.com/jeffscottbrown/applemusic/server"
)

func main() {
	logging.Configure()

	server.RunServer()
}
