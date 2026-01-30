package main

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

type liveReloadInjector struct {
	http.ResponseWriter
	injected bool
}

func (w *liveReloadInjector) WriteHeader(statusCode int) {
	if strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		w.Header().Del("Content-Length")
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *liveReloadInjector) Write(b []byte) (int, error) {
	if !w.injected && strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		w.Header().Del("Content-Length")

		content := string(b)
		if idx := strings.LastIndex(content, "</body>"); idx != -1 {
			content = content[:idx] + reloaderScript + content[idx:]
			w.injected = true
			return w.ResponseWriter.Write([]byte(content))
		}
	}
	return w.ResponseWriter.Write(b)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	log.Printf("Browser connected (total: %d)", len(clients))

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
	log.Printf("Browser disconnected (total: %d)", len(clients))
}

func notifyClients() {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
			log.Printf("Error notifying client: %v", err)
			conn.Close()
			delete(clients, conn)
		}
	}
	log.Printf("Notified %d browser(s) to reload", len(clients))
}
