package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/system"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	clients   = make(map[*websocket.Conn]*models.ClientInfo)
	clientsMu sync.Mutex
)

const debug = false

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	clientInfo := &models.ClientInfo{
		ID:          "",
		ConnectedAt: time.Now(),
	}

	clientsMu.Lock()
	clients[conn] = clientInfo
	clientsMu.Unlock()

	log.Println("client connected")

	go func(c *websocket.Conn) {
		defer func() {
			clientsMu.Lock()
			delete(clients, c)
			clientsMu.Unlock()
			c.Close()
			log.Println("client disconnected")
		}()

		c.SetReadDeadline(time.Now().Add(config.AppConfig.PongWait))
		c.SetPongHandler(func(string) error {
			c.SetReadDeadline(time.Now().Add(config.AppConfig.PongWait))
			return nil
		})

		pingTicker := time.NewTicker(config.AppConfig.PingPeriod)
		defer pingTicker.Stop()

		go func() {
			for range pingTicker.C {
				if err := c.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(config.AppConfig.WriteDeadline)); err != nil {
					return
				}
			}
		}()

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				break
			}

			var msgData map[string]interface{}
			if err := json.Unmarshal(message, &msgData); err == nil {
				if id, ok := msgData["id"].(string); ok && id != "" {
					clientsMu.Lock()
					if info, exists := clients[c]; exists {
						info.ID = id
					}
					clientsMu.Unlock()
				}
			}

			if debug {
				log.Println(string(message))
			}
		}
	}(conn)
}

func BroadcastLoop(ctx context.Context) {
	ticker := time.NewTicker(config.AppConfig.BroadcastRate)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cursor := system.GetCursorPosition()

			windows, err := system.GetCachedChromeWindows()
			if err != nil {
				log.Println("getCachedChromeWindows error:", err)
				continue
			}

			positions := []models.MessageInfo{}
			for _, window := range windows {
				positions = append(positions, models.MessageInfo{
					ID:    window.Title,
					X:     cursor.X - window.X,
					Y:     cursor.Y - window.Y - config.AppConfig.ChromeOffsetY,
					GridX: -window.X,
					GridY: -window.Y - config.AppConfig.ChromeOffsetY,
				})
			}

			data, err := json.Marshal(positions)
			if err != nil {
				log.Println("json marshal error:", err)
				continue
			}

			clientsMu.Lock()
			for c := range clients {
				c.SetWriteDeadline(time.Now().Add(config.AppConfig.WriteDeadline))
				if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Println("write error, removing client:", err)
					c.Close()
					delete(clients, c)
				}
			}
			clientsMu.Unlock()
		}
	}
}
