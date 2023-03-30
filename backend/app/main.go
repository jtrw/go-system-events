package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"
    "github.com/gorilla/websocket"
)

// Manager is used to hold references to all Clients Registered, and Broadcasting etc
type Manager struct {
    clients ClientList

    // Using a syncMutex here to be able to lcok state before editing clients
    // Could also use Channels to block
    sync.RWMutex
}

// NewManager is used to initalize all the values inside the manager
func NewManager() *Manager {
	return &Manager{
	    clients: make(ClientList),
	}
}

var websocketUpgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type ClientList map[*Client]bool

type Client struct {
	// the websocket connection
	connection *websocket.Conn

	// manager is the manager used to manage the client
	manager *Manager
	// egress is used to avoid concurrent writes on the WebSocket
    egress chan []byte
}

var event string

// serveWS is a HTTP Handler that the has the Manager that allows connections
func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {

	log.Println("New connection")

	websocketUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// Begin by upgrading the HTTP request
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Create New Client
    client := NewClient(conn, m)
    m.addClient(client)
    // Start the read / write processes
    go client.readMessages()
    // go client.writeMessages()
	// We wont do anything yet so close connection again
	//conn.Close()
}

func (c *Client) readMessages() {
	defer func() {
		// Graceful Close the Connection once this
		// function is done
		c.manager.removeClient(c)
	}()
	// Loop Forever
	for {
		// ReadMessage is used to read the next message in queue
		// in the connection
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but not simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break // Break the loop to close conn & Cleanup
		}
		log.Println("MessageType: ", messageType)
		log.Println("Payload: ", string(payload))

		// Hack to test that WriteMessages works as intended
        // Will be replaced soon
        for wsclient := range c.manager.clients {
            wsclient.egress <- payload
        }
	}
}

// writeMessages is a process that listens for new messages to output to the Client
func (c *Client) writeMessages() {
	defer func() {
		// Graceful close if this triggers a closing
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			// Ok will be false Incase the egress channel is closed
			if !ok {
				// Manager has closed this connection channel, so communicate that to frontend
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					// Log that the connection is closed and the reason
					log.Println("connection closed: ", err)
				}
				// Return to close the goroutine
				return
			}
			// Write a Regular text message to the connection
			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println(err)
			}
			log.Println("sent message")
		}
	}
}

// addClient will add clients to our clientList
func (m *Manager) addClient(client *Client) {
	// Lock so we can manipulate
	m.Lock()
	defer m.Unlock()

	// Add Client
	m.clients[client] = true
}

// removeClient will remove the client and clean up
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// Check if Client exists, then delete it
	if _, ok := m.clients[client]; ok {
		// close connection
		client.connection.Close()
		// remove
		delete(m.clients, client)
	}
}

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Home Page")
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
