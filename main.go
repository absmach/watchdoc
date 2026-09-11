// Package main provides WatchDoc, a live-reload development server.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/skratchdot/open-golang/open"
)

type options struct {
	port      string
	serveDir  string
	watchDirs string
	cmd       string
	noBrowser bool
}

func parseOptions(args []string, output io.Writer) (options, error) {
	var opts options
	flags := flag.NewFlagSet("watchdoc", flag.ContinueOnError)
	flags.SetOutput(output)
	flags.StringVar(&opts.port, "port", "8080", "port to run the file server on")
	flags.StringVar(&opts.serveDir, "serve-dir", ".", "directory to serve files from")
	flags.StringVar(&opts.watchDirs, "watch-dirs", "", "additional comma-separated directories to watch")
	flags.StringVar(&opts.cmd, "cmd", "", "command to execute on file change")
	flags.BoolVar(&opts.noBrowser, "no-browser", false, "disable automatic browser opening")

	if err := flags.Parse(args); err != nil {
		return options{}, err
	}
	if flags.NArg() != 0 {
		return options{}, fmt.Errorf(
			"unexpected positional arguments: %q; boolean flags take no spaced value (use -no-browser or -no-browser=true)",
			strings.Join(flags.Args(), " "),
		)
	}

	return opts, nil
}

func main() {
	opts, err := parseOptions(os.Args[1:], os.Stderr)
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	absPath, err := filepath.Abs(opts.serveDir)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("File server started at port %s", opts.port)
	log.Printf("Open your browser at http://localhost:%s", opts.port)
	if opts.cmd != "" {
		log.Printf("Command: %s", opts.cmd)
	}
	if opts.watchDirs != "" {
		log.Printf("Watching directories: %s", strings.Join([]string{absPath, opts.watchDirs}, ","))
	}
	log.Printf("Serving from: %s", absPath)

	serveList, watchList := resolveWatchDirs(opts.watchDirs, absPath)
	go watchFiles(serveList, watchList, opts.cmd, absPath)

	http.HandleFunc("/ws", handleWebSocket)

	fileServer := http.FileServer(http.Dir(opts.serveDir))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		injector := &liveReloadInjector{ResponseWriter: w}
		fileServer.ServeHTTP(injector, r)
	})

	srv := &http.Server{
		Addr:              ":" + opts.port,
		ReadHeaderTimeout: 10 * time.Second,
	}

	if !opts.noBrowser {
		go func() {
			if err := open.Start("http://localhost:" + opts.port); err != nil {
				log.Printf("Failed to open browser: %v", err)
			}
		}()
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
