// Package ws is the structured websocket event hub. The simulator broadcasts
// JSON events (cloud.changed/state.changed/scenario.loaded/...) to every
// connected client.
package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Event is the structured message broadcast to all websocket clients.
type Event struct {
	Type         string `json:"type"`                   // cloud.changed | state.changed | scenario.loaded | reset | hello
	Scope        string `json:"scope"`                  // cloud | state | sim
	ResourceType string `json:"resourceType,omitempty"` // e.g. nimbus_subnet
	ID           string `json:"id,omitempty"`           // cloud id
	Serial       int    `json:"serial,omitempty"`
	Detail       string `json:"detail,omitempty"`
	Timestamp    string `json:"timestamp"`
}

// Broadcast marshals and sends a structured event to all connected clients,
// stamping the timestamp if unset.
func Broadcast(ev Event) {
	if SOCKEX == nil {
		return
	}
	if ev.Timestamp == "" {
		ev.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return
	}
	SOCKEX.send(b)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Init creates the global exchange if it does not already exist.
func Init() {
	if SOCKEX == nil {
		SOCKEX = newSockEx()
	}
}

// SocketHandler upgrades a connection and registers it with the exchange.
func SocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	Init()
	SOCKEX.add(conn)
}

// SOCKEX is the global socket exchange instance.
var SOCKEX *socketExchange

// client wraps a connection with a buffered outbound queue. A single writePump
// goroutine drains the queue, so writes to one connection are never concurrent
// (gorilla/websocket forbids that).
type client struct {
	conn *websocket.Conn
	send chan []byte
}

// socketExchange tracks connected clients and fans messages out to them.
type socketExchange struct {
	mu      sync.Mutex
	lastID  int
	clients map[int]*client
}

func newSockEx() *socketExchange {
	return &socketExchange{clients: map[int]*client{}}
}

// add registers a new connection, starts its read/write pumps, and greets it.
func (se *socketExchange) add(c *websocket.Conn) {
	cl := &client{conn: c, send: make(chan []byte, 32)}
	se.mu.Lock()
	se.lastID++
	id := se.lastID
	se.clients[id] = cl
	se.mu.Unlock()

	go se.writePump(cl)
	go se.readPump(id, cl)

	hello, _ := json.Marshal(Event{Type: "hello", Scope: "sim", Timestamp: time.Now().UTC().Format(time.RFC3339)})
	cl.enqueue(hello)
}

// send fans a message out to every client's queue (non-blocking; a client whose
// queue is full is dropped from this message rather than stalling the broadcast).
func (se *socketExchange) send(msg []byte) {
	se.mu.Lock()
	clients := make([]*client, 0, len(se.clients))
	for _, cl := range se.clients {
		clients = append(clients, cl)
	}
	se.mu.Unlock()
	for _, cl := range clients {
		cl.enqueue(msg)
	}
}

func (cl *client) enqueue(msg []byte) {
	defer func() { _ = recover() }() // tolerate send on a closed channel during teardown
	select {
	case cl.send <- msg:
	default: // queue full: drop this message for this slow client
	}
}

func (se *socketExchange) remove(id int) {
	se.mu.Lock()
	cl, ok := se.clients[id]
	if ok {
		delete(se.clients, id)
		close(cl.send)
	}
	se.mu.Unlock()
}

// writePump is the sole writer for a connection.
func (se *socketExchange) writePump(cl *client) {
	for msg := range cl.send {
		if err := cl.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
	_ = cl.conn.Close()
}

// readPump drains inbound frames (we don't act on them) and tears the client down
// on disconnect.
func (se *socketExchange) readPump(id int, cl *client) {
	defer se.remove(id)
	for {
		if _, _, err := cl.conn.NextReader(); err != nil {
			_ = cl.conn.Close()
			return
		}
	}
}
