package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestLiveReloadInjector_HTMLWithBody(t *testing.T) {
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "text/html; charset=utf-8")

	injector := &liveReloadInjector{ResponseWriter: rec}
	input := []byte("<html><body><p>Hello</p></body></html>")

	n, err := injector.Write(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == 0 {
		t.Fatal("expected bytes written > 0")
	}

	body := rec.Body.String()
	if !strings.Contains(body, reloaderScript) {
		t.Error("expected reload script to be injected")
	}
	if !strings.Contains(body, "</body>") {
		t.Error("expected </body> tag to be preserved")
	}
	if !injector.injected {
		t.Error("expected injected flag to be true")
	}
}

func TestLiveReloadInjector_HTMLWithoutBody(t *testing.T) {
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "text/html; charset=utf-8")

	injector := &liveReloadInjector{ResponseWriter: rec}
	input := []byte("<html><p>No body tag</p></html>")

	_, err := injector.Write(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := rec.Body.String()
	if strings.Contains(body, reloaderScript) {
		t.Error("should not inject script when no </body> tag")
	}
}

func TestLiveReloadInjector_NonHTML(t *testing.T) {
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "application/json")

	injector := &liveReloadInjector{ResponseWriter: rec}
	input := []byte(`{"key": "value"}`)

	_, err := injector.Write(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := rec.Body.String()
	if strings.Contains(body, "<script>") {
		t.Error("should not inject script into non-HTML response")
	}
	if body != `{"key": "value"}` {
		t.Errorf("body should be unchanged, got: %s", body)
	}
}

func TestLiveReloadInjector_WriteHeader_RemovesContentLength(t *testing.T) {
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "text/html")
	rec.Header().Set("Content-Length", "100")

	injector := &liveReloadInjector{ResponseWriter: rec}
	injector.WriteHeader(http.StatusOK)

	if cl := rec.Header().Get("Content-Length"); cl != "" {
		t.Errorf("Content-Length should be removed for HTML, got: %s", cl)
	}
}

func TestLiveReloadInjector_WriteHeader_KeepsContentLengthForNonHTML(t *testing.T) {
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "application/json")
	rec.Header().Set("Content-Length", "100")

	injector := &liveReloadInjector{ResponseWriter: rec}
	injector.WriteHeader(http.StatusOK)

	if cl := rec.Header().Get("Content-Length"); cl != "100" {
		t.Errorf("Content-Length should be preserved for non-HTML, got: %s", cl)
	}
}

func TestHandleWebSocket(t *testing.T) {
	// Reset global state
	clientsMu.Lock()
	clients = make(map[*websocket.Conn]bool)
	clientsMu.Unlock()

	srv := httptest.NewServer(http.HandlerFunc(handleWebSocket))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	// Give server time to register client
	time.Sleep(50 * time.Millisecond)

	clientsMu.Lock()
	count := len(clients)
	clientsMu.Unlock()

	if count != 1 {
		t.Errorf("expected 1 client, got %d", count)
	}

	conn.Close()

	// Give server time to deregister
	time.Sleep(50 * time.Millisecond)

	clientsMu.Lock()
	count = len(clients)
	clientsMu.Unlock()

	if count != 0 {
		t.Errorf("expected 0 clients after disconnect, got %d", count)
	}
}

func TestNotifyClients(t *testing.T) {
	clientsMu.Lock()
	clients = make(map[*websocket.Conn]bool)
	clientsMu.Unlock()

	srv := httptest.NewServer(http.HandlerFunc(handleWebSocket))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	time.Sleep(50 * time.Millisecond)

	notifyClients()

	conn.SetReadDeadline(time.Now().Add(time.Second))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}
	if string(msg) != "reload" {
		t.Errorf("expected 'reload', got '%s'", string(msg))
	}
}
