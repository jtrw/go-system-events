package websocket

import (
    "log"
    "time"
    "encoding/json"
    "github.com/gorilla/websocket"
)

// ClientList is a map used to help manage a map of clients
type ClientList map[*Client]bool

// Client is a websocket client, basically a frontend visitor
type Client struct {
	// the websocket connection
	connection *websocket.Conn
	// manager is the manager used to manage the client
	manager *Manager
    egress chan Event
}

var (
	// pongWait is how long we will await a pong response from client
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10
)

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
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
        _, payload, err := c.connection.ReadMessage()

        if err != nil {
            // If Connection is closed, we will Recieve an error here
            // We only want to log Strange errors, but simple Disconnection
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error reading message: %v", err)
            }
            break // Break the loop to close conn & Cleanup
        }
        // Marshal incoming data into a Event struct
        var request Event
        if err := json.Unmarshal(payload, &request); err != nil {
            log.Printf("error marshalling message: %v", err)
            break // Breaking the connection here might be harsh xD
        }
        log.Println(string(request.Type))
        log.Println(string(request.Token))
        log.Println(string(request.Payload))
        // Route the Event
        if err := c.manager.routeEvent(request, c); err != nil {
            log.Println("Error handeling Message: ", err)
        }
	}
}

// pongHandler is used to handle PongMessages for the Client
func (c *Client) pongHandler(pongMsg string) error {
	// Current time + Pong Wait time
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

// writeMessages is a process that listens for new messages to output to the Client
func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingInterval)
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

                data, err := json.Marshal(message)
                if err != nil {
                    log.Println(err)
                    return // closes the connection, should we really
                }
                // Write a Regular text message to the connection
                if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
                    log.Println(err)
                }
                log.Println("sent message")
            case <-ticker.C:
                log.Println("ping")
                // Send the Ping
                if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
                    log.Println("writemsg: ", err)
                    return // return to break this goroutine triggeing cleanup
                }
        }
	}
}
