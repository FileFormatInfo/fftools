package internal

import (
	"fmt"
	"log/slog"
)

var (
	BUILDER = "unknown"
	COMMIT  = "(local)"
	LASTMOD = "(local)"
	VERSION = "internal"
)

func PrintVersion(name string) {
	if LogLevel >= slog.LevelInfo {
		slog.Info("Version information", "name", name, "version", VERSION, "lastmod", LASTMOD, "commit", COMMIT, "builder", BUILDER)
	} else {
		fmt.Printf("%s version %s (built on %s from %s by %s)\n", name, VERSION, LASTMOD, COMMIT, BUILDER)
	}
}
