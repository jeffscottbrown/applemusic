package logging

import "log/slog"

func Configure() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}
