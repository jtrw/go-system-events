package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

// Manager is used to hold references to all Clients Registered, and Broadcasting etc
type Manager struct {
}

// NewManager is used to initalize all the values inside the manager
func NewManager() *Manager {
	return &Manager{}
}

var websocketUpgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

var event string

// serveWS is a HTTP Handler that the has the Manager that allows connections
func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {

	log.Println("New connection")
	// Begin by upgrading the HTTP request
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// We wont do anything yet so close connection again
	conn.Close()
}

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Home Page")
}



func reader(conn *websocket.Conn) {
    for {
    // read in a message
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
    // print out that message for clarity
        fmt.Println(string(p))

        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println(err)
            return
        }
    }
}


func setupRoutes() {
    manager := NewManager()

    http.HandleFunc("/", homePage)
    http.HandleFunc("/ws", manager.serveWS)
}

func main() {
    fmt.Println("Hello World")
    setupRoutes()
    log.Fatal(http.ListenAndServe(":8080", nil))
}
