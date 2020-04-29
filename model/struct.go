package model

import "github.com/gorilla/websocket"

// Files are saved as base64 format to database then if queried,
const (
	TXT = iota
	MP3
	EXE
	MP4
	WAV
	JPG
	PNG
)

type CookieDetail struct {
	Email      string
	Collection string
	CookieName string
	Path       string
	Data       map[string]interface{}
}

type User struct {
	Email string `bson:"_id" json:"email"`
	Name  string `bson:"name" json:"name"`
	DOB   string `bson:"age" json:"age"`
	Class string `bson:"class" json:"class"`
	// ID should either users matric or email stripping @....
	ID           string   `bson:"userID" json:"userID"`
	ParentEmail  string   `bson:"parentEmail" json:"parentEmail"`
	ParentNumber string   `bson:"parentNumber" json:"parentNumber"`
	Password     []byte   `bson:"password" json:"password"`
	Faculty      string   `bson:"faculty" json:"faculty"`
	UUID         string   `bson:"loginUUID" json:"uuid"`
	RoomsJoined  []string `bson:"roooms" json:"rooms"`
}

type Admin struct {
	StaffDetails User `bson:",inline"`
	Super        bool `bson:"super" json:"super"`
}

type Chats struct {
	RoomID          string    `bson:"_id" json:"email"`
	RegisteredUsers []string  `bson:"registeredUsers"`
	Messages        []Message `bson:"messages" json:"messages"`
}

type Message struct {
	RoomID  string `bson:"omitempty" json:"roomID"`
	Message string `bson:"message" json:"message"`
	User    string `bson:"user" json:"user"`
	Index   int    `bson:"index" json:"index"`
	Type    int    `bson:"type" json:"type"`
}

type Joined struct {
	RoomID string `json:"roomID"`
	Joined bool   `json:"joined"`
}

// FileType save files separately and make sure they are distinct
type FileType struct {
	Downloaded bool   `bson:"downloaded" json:"downloaded"`
	Sha256     string `bson:"_id" json:"sha256"`
}

type WSMessage struct {
	Data []byte
	User string
}

type Subscription struct {
	Conn *Connection
	User string
}

// connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The websocket connection.
	WS *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	Users map[string]map[*Connection]bool

	// Inbound messages from the connections.
	Broadcast chan WSMessage

	// Register requests from the connections.
	Register chan Subscription

	// Unregister requests from connections.
	UnRegister chan Subscription
}
