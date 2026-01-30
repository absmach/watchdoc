package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func main() {
	p := flag.String("port", "8080", "port to run the file server on")
	dir := flag.String("dir", "", "directory to serve files from (default: smart detection)")
	extraDirs := flag.String("watch-dirs", "", "additional comma-separated directories to watch")
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
	if *extraDirs != "" {
		log.Printf("Watching directories: %s", strings.Join([]string{absPath, *extraDirs}, ","))
	}
	log.Printf("Serving from: %s", absPath)

	watchList, extrasList := resolveWatchDirs(*extraDirs, absPath)
	go watchFiles(watchList, extrasList, *cmd, absPath)

	http.HandleFunc("/ws", handleWebSocket)

	fileServer := http.FileServer(http.Dir(serveDir))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		injector := &liveReloadInjector{ResponseWriter: w}
		fileServer.ServeHTTP(injector, r)
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
