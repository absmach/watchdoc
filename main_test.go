package main

import (
	"io"
	"strings"
	"testing"
)

func TestParseOptions(t *testing.T) {
	t.Run("supported no-browser form preserves later flags", func(t *testing.T) {
		opts, err := parseOptions([]string{
			"-no-browser",
			"-watch-dirs", "./content",
			"-serve-dir", "./out",
			"-cmd", "./build.bat",
		}, io.Discard)
		if err != nil {
			t.Fatalf("parseOptions() error = %v", err)
		}
		if !opts.noBrowser || opts.watchDirs != "./content" || opts.serveDir != "./out" || opts.cmd != "./build.bat" {
			t.Fatalf("parseOptions() = %+v", opts)
		}
	})

	t.Run("spaced boolean value is rejected", func(t *testing.T) {
		_, err := parseOptions([]string{
			"-no-browser", "true",
			"-watch-dirs", "./content",
			"-serve-dir", "./out",
			"-cmd", "./build.bat",
		}, io.Discard)
		if err == nil {
			t.Fatal("parseOptions() error = nil")
		}
		if !strings.Contains(err.Error(), "-no-browser=true") {
			t.Fatalf("parseOptions() error = %q, want boolean flag hint", err)
		}
	})
}
