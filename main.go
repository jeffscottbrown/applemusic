package main

import (
	"github.com/jeffscottbrown/applemusic/commit"
	"github.com/jeffscottbrown/applemusic/logging"
	"github.com/jeffscottbrown/applemusic/server"
	"log/slog"
)

func main() {
	logging.Configure()

	slog.Debug("Build Info", "Build Time", commit.BuildTime, "Commit Tag", commit.Tag)

	server.Run()
}
