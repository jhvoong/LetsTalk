package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"

	"github.com/metaclips/LetsTalk/backend/values"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 30 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 120 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var HubConstruct = Hub{
	Broadcast:  make(chan WSMessage),
	Register:   make(chan Subscription),
	UnRegister: make(chan Subscription),
	Users:      make(map[string]map[*Connection]bool),
}

func (h *Hub) Run() {
	Upgrader.CheckOrigin = func(r *http.Request) bool {
		// ToDo: save host in config.json or .env
		host := r.Header.Get("Origin")
		if host == "https://127.0.0.1:8081" && r.TLS != nil && r.TLS.HandshakeComplete {
			return true
		}

		return false
	}

	for {
		select {
		case s := <-h.Register:
			connections := h.Users[s.User]
			if connections == nil {
				connections = make(map[*Connection]bool)
				h.Users[s.User] = connections
			}
			h.Users[s.User][s.Conn] = true
			fmt.Println("User registered")
			go broadcastOnlineStatusToAllUserRoom(s.User, true)

		case s := <-h.UnRegister:
			connections := h.Users[s.User]
			if connections != nil {
				if _, ok := connections[s.Conn]; ok {
					delete(connections, s.Conn)
					close(s.Conn.Send)
					if len(connections) == 0 {
						delete(h.Users, s.User)
						go broadcastOnlineStatusToAllUserRoom(s.User, false)
						fmt.Println(s.User, "offline")
					}
					fmt.Println(s.User, "subscription removed")
				}
			}

		case m := <-h.Broadcast:
			connections := h.Users[m.User]
			for c := range connections {
				select {
				case c.Send <- m.Data:
				default:
					close(c.Send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.Users, m.User)
						go broadcastOnlineStatusToAllUserRoom(m.User, false)
					}
				}
			}
		}
	}
}

func (h *Hub) sendMessage(msg []byte, user string) {
	m := WSMessage{msg, user}
	HubConstruct.Broadcast <- m
}

// WritePump pumps messages from the hub to the websocket connection.
func (s *Subscription) WritePump() {
	c := s.Conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.WS.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	if err := c.WS.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}

	return c.WS.WriteMessage(mt, payload)
}

// ReadPump pumps messages from the websocket connection to the hub.
func (s Subscription) ReadPump(user string) {
	c := s.Conn

	defer func() {
		HubConstruct.UnRegister <- s
		c.WS.Close()
	}()

	c.WS.SetReadDeadline(time.Now().Add(pongWait))

	c.WS.SetPongHandler(
		func(string) error {
			return c.WS.SetReadDeadline(time.Now().Add(pongWait))
		})

	for {
		var err error
		var msg messageBytes
		_, msg, err = c.WS.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v\n", err)
			}
			log.Println(err)

			break
		}

		// TODO: Always enquire for userID.
		data := struct {
			MsgType    string `json:"msgType"`
			User       string `json:"userID"`
			SearchText string `json:"searchText"`
		}{}

		err = json.Unmarshal(msg, &data)
		if err != nil {
			log.Println("could not unmarshal json")
		}

		if data.User != user {
			log.Warningf("an invalidated user tried to make websocket request for %s, invalidated user: %s, user: %s", data.MsgType, data.User, user)
			continue
		}

		switch data.MsgType {
		// TODO: add support to remove message.
		// TODO: users should choose if to join chat.
		case values.WebsocketOpenMsgType:
			handleLoadUserContent(user)

		case values.RequestMessages:
			msg.handleRequestMessages(user)

		case values.NewRoomCreatedMsgType:
			msg.handleCreateNewRoom()

		case values.ExitRoomMsgType:
			msg.handleExitRoom(user)

		case values.RequestUsersToJoinRoomMsgType:
			msg.handleRequestUserToJoinRoom()

		case values.UserJoinedRoomMsgType:
			msg.handleUserAcceptRoomRequest(user)

		case values.NewFileUploadMsgType:
			msg.handleNewFileUpload()

		case values.UploadFileChunkMsgType:
			msg.handleUploadFileChunk()

		case values.UploadFileSuccessMsgType:
			msg.handleUploadFileUploadComplete()

		case values.RequestDownloadMsgType:
			msg.handleRequestDownload(user)

		case values.DownloadFileChunkMsgType:
			msg.handleFileDownload(user)

		case values.StartClassSession:
			classSessions.startClassSession(msg, user)

		case values.JoinClassSession:
			classSessions.joinClassSession(msg, user)

		case values.NegotiateSDP:
			sdpConstruct{}.acceptRenegotiation(msg)

		case values.NewMessageMsgType:
			msg.handleNewMessage(user)

		case values.SearchUserMsgType:
			handleSearchUser(data.SearchText, user)

		default:
			log.Println("Could not convert required type", data.MsgType)
		}
	}
}
