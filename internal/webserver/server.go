package webserver

import (
	"fmt"
	"net"
	"net/http"

	"github.com/leetrout/terraform-sim/internal/cloud"
	"github.com/leetrout/terraform-sim/internal/sim"
	"github.com/leetrout/terraform-sim/internal/ws"
)

// handler serves the embedded SPA index for any non-API, non-static path.
func handler(w http.ResponseWriter, r *http.Request) {
	html, err := f.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "index not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(html)
}

// NewMux builds the application's HTTP routes: the cloud + sim APIs, the websocket
// endpoint, and the embedded frontend.
func NewMux(cloudStore *cloud.Store, simulator *sim.Simulator) *http.ServeMux {
	mux := http.NewServeMux()
	cloud.NewHandler(cloudStore).Register(mux)
	sim.NewHandler(simulator).Register(mux)
	mux.HandleFunc("/ws", ws.SocketHandler)
	mux.Handle("/static/", http.FileServer(http.FS(f)))
	mux.HandleFunc("/", handler)
	return mux
}

// Serve binds the listener first (so the browser is only opened once the server
// is actually accepting connections) and then serves. onListen is invoked with
// the bound address after the listener is up.
func Serve(addr string, mux *http.ServeMux, onListen func(addr string)) error {
	ws.Init()
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}
	if onListen != nil {
		onListen(ln.Addr().String())
	}
	return http.Serve(ln, mux)
}
