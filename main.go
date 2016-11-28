package main

import (
	"fmt"
	"net/http"
	"time"

	"reflect"

	"golang.org/x/net/websocket"
)

type hub struct {
	clients    map[*client]bool
	broadcast  chan string
	register   chan *client
	unregister chan *client

	content string
}

var h = hub{
	broadcast:  make(chan string),
	register:   make(chan *client),
	unregister: make(chan *client),
	clients:    make(map[*client]bool),
	content:    "",
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			fmt.Println("registering client")
			h.clients[c] = true
			c.send <- []byte(h.content)
			break

		case c := <-h.unregister:
			fmt.Println("unregistering client")

			_, ok := h.clients[c]
			if ok {
				delete(h.clients, c)
				close(c.send)
			}
			break

		case m := <-h.broadcast:
			fmt.Println("broadcast")
			h.content = m
			h.broadcastMessage()
			break
		}
	}
}

func (h *hub) broadcastMessage() {
	for c := range h.clients {
		select {
		case c.send <- []byte(h.content):
			break

		// We can't reach the client
		default:
			close(c.send)
			delete(h.clients, c)
		}
	}
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

type client struct {
	ws    *websocket.Conn
	send  chan []byte
	isReg bool
	name  string
}

func (c *client) writePump() {
	fmt.Println("Write Pump")
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		fmt.Println("Closing from write")
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			fmt.Println(reflect.TypeOf(message))
			if !ok {
				c.write([]byte{})
				return
			}
			if err := c.write(message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write([]byte{}); err != nil {
				return
			}
		}
	}
}

func (c *client) write(message []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	_, err := c.ws.Write([]byte(message))
	return err
}

func (c *client) readPump() {
	fmt.Println("Read Pump")

	defer func() {
		h.unregister <- c
		fmt.Println("Closing from read")
		c.ws.Close()
	}()

	c.ws.SetReadDeadline(time.Now().Add(pongWait))

	for {
		var msg string
		for {
			err := websocket.Message.Receive(c.ws, &msg)
			if err != nil {
				break
			}
			if !c.isReg {
				c.name = msg
				c.isReg = true
				continue
			}
			fmt.Printf("Received: %s.\n", msg)
			message := c.name + ": " + msg
			h.broadcast <- message
			fmt.Println("current clients: ", h.clients)
		}
	}
}

func serveWs(ws *websocket.Conn) {
	c := &client{
		send:  make(chan []byte, maxMessageSize),
		ws:    ws,
		isReg: false,
		name:  "",
	}
	websocket.Message.Send(c.ws, "What name would you like")
	fmt.Println("Client connected")
	fmt.Println("register client")
	h.register <- c
	fmt.Println("start write pump")

	go c.writePump()
	fmt.Println("start read pump")

	c.readPump()
}

func main() {
	go h.run()
	http.Handle("/ws", websocket.Handler(serveWs))
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
