package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// Client represents a single SSE connection
type Client struct {
	channel chan string
}

// clients holds all active SSE connections
var clients = make(map[Client]bool)

// sendMessage sends a message to all connected SSE clients
func SendMessage(message string) {
	for client := range clients {
		client.channel <- message
	}
}

// sseHandler handles new client connections for SSE
func RegisterSSEHandlers(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("events", func(c echo.Context) error {
			client := Client{channel: make(chan string)}
			clients[client] = true

			c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
			c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
			c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
			c.Response().Flush()

			// When the client closes the connection, remove them from the clients map
			defer func() {
				delete(clients, client)
				close(client.channel)
			}()

			// Start an endless loop, sending messages as they come
			for {
				message, open := <-client.channel
				if !open {
					break
				}

				if _, err := c.Response().Write([]byte("data: " + message + "\n\n")); err != nil {
					c.Echo().Logger.Error(err)
					break
				}
				c.Response().Flush()
			}

			return nil
		})
		return nil
	})
}
