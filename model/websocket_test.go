package model

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/metaclips/LetsTalk/values"
	"go.uber.org/goleak"
)

func TestWS(t *testing.T) {
	defer func() {
		<-ctx.Done()
		goleak.VerifyNone(t)
	}()

	if err := values.LoadConfiguration("../config.json"); err != nil {
		t.Fatal("could not load configurations", err)
	}

	addDummyUser()
	InitModel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := Upgrader.Upgrade(w, r, nil)
		t.Error(err)
		defer ws.Close()

		conn := &Connection{Send: make(chan []byte, 256), WS: ws}

		s := Subscription{Conn: conn, User: "Test@email.com"}
		HubConstruct.Register <- s
		t.Log("User test connected")

		go s.ReadPump("Test@email.com")
		s.WritePump()
		t.Log("Websocket closed")
		cancel()
	}))

	u := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the server.
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("here %v", err)
	}
	defer ws.Close()

	for i := 0; i < 10; i++ {
		if err := ws.WriteMessage(websocket.TextMessage, []byte("Hello")); err != nil {
			t.Fatalf("here %v", err)
		}

		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("here %v", err)
		}

		if string(p) != "Hi" {
			t.Fatalf("bad message")
		}
	}
}

func addDummyUser() error {
	// Password is generated from name.
	user := User{
		Name:  "Test",
		Email: "Test@email.com",
	}

	err := user.UploadUser()
	return err
}

func createNewRoom() {
	jsonStruct := struct {
		Email   string `json:"email"`
		MsgType string `json:"msgType"`
		RoomName string `json:"roomName"`
	}{
		"Test@email.com"
	}
}
