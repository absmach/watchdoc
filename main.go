// Package main provides WatchDoc, a live-reload development server.
package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	p := flag.String("port", "8080", "port to run the file server on")
	dir := flag.String("serve-dir", "", "directory to serve files from (default: current directory)")
	watchDirs := flag.String("watch-dirs", "", "additional comma-separated directories to watch")
	cmd := flag.String("cmd", "", "command to execute on file change")
	flag.Parse()

	port := *p
	serveDir := *dir

	if serveDir == "" {
		serveDir = "."
	}

	absPath, err := filepath.Abs(serveDir)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("File server started at port %s", port)
	log.Printf("Open your browser at http://localhost:%s", port)
	if *cmd != "" {
		log.Printf("Command: %s", *cmd)
	}
	if *watchDirs != "" {
		log.Printf("Watching directories: %s", strings.Join([]string{absPath, *watchDirs}, ","))
	}
	log.Printf("Serving from: %s", absPath)

	watchList, extrasList := resolveWatchDirs(*watchDirs, absPath)
	go watchFiles(watchList, extrasList, *cmd, absPath)

	http.HandleFunc("/ws", handleWebSocket)

	fileServer := http.FileServer(http.Dir(serveDir))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		injector := &liveReloadInjector{ResponseWriter: w}
		fileServer.ServeHTTP(injector, r)
	})

	srv := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 10 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
