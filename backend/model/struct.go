package model

import (
	"os"
	"sync"
	"time"

	"github.com/at-wat/ebml-go/webm"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/sharing"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media/samplebuilder"
)

type CookieDetail struct {
	Email      string
	Collection string
	CookieName string
	Path       string
	Data       CookieData
}

type CookieData struct {
	ExitTime time.Time
	UUID     string
	Email    string
	Super    bool
}

type User struct {
	IsAdmin bool   `bson:"isAdmin"`
	Email   string `bson:"_id" json:"email"`
	Name    string `bson:"name" json:"name"`
	DOB     string `bson:"age" json:"age"`
	Class   string `bson:"class" json:"class"`
	// ID should either be users matric or leading email stripping @....
	ID           string        `bson:"userID" json:"userID"`
	ParentEmail  string        `bson:"parentEmail" json:"parentEmail"`
	ParentNumber string        `bson:"parentNumber" json:"parentNumber"`
	Password     []byte        `bson:"password" json:"password"`
	Faculty      string        `bson:"faculty" json:"faculty"`
	UUID         string        `bson:"loginUUID" json:"uuid"`
	RoomsJoined  []RoomsJoined `bson:"roomsJoined" json:"roomsJoined"`
	JoinRequest  []JoinRequest `bson:"joinRequest" json:"joinRequest"`
}

type RoomsJoined struct {
	RoomID   string `bson:"rooomID" json:"roomID"`
	RoomName string `bson:"rooomName" json:"roomName"`
}

type JoinRequest struct {
	RoomID             string   `bson:"_id" json:"roomID"`
	RoomName           string   `bson:"roomName" json:"roomName"`
	RequestingUserName string   `bson:"requestingUserName" json:"requestingUserName"`
	RequestingUserID   string   `bson:"requestingUserID" json:"requestingUserID"`
	Users              []string `bson:"-" json:"users"`
}

type Admin struct {
	StaffDetails User `bson:",inline"`
	Super        bool `bson:"super" json:"super"`
}

type Room struct {
	RoomID          string    `bson:"_id" json:"email"`
	RoomName        string    `bson:"roomName" json:"roomName"`
	RegisteredUsers []string  `bson:"registeredUsers"`
	Messages        []Message `bson:"messages" json:"messages"`
}

type Message struct {
	RoomID      string `bson:"-" json:"roomID,omitempty"`
	Message     string `bson:"message" json:"message,omitempty"`
	FileSize    string `bson:"fileSize,omitempty" json:"fileSize,omitempty"`
	FileHash    string `bson:"fileHash,omitempty" json:"fileHash,omitempty"`
	UserID      string `bson:"userID" json:"userID,omitempty"`
	Name        string `bson:"name" json:"name,omitempty"`
	Index       int    `bson:"index" json:"index,omitempty"`
	Time        string `bson:"time" json:"time,omitempty"`
	Type        string `bson:"type" json:"type,omitempty"`
	MessageType string `bson:"-" json:"msgType,omitempty"`
}

type Joined struct {
	RoomID      string `json:"roomID"`
	RoomName    string `json:"roomName"`
	Email       string `json:"email"`
	Joined      bool   `json:"joined"`
	MessageType string `bson:"-" json:"msgType"`
}

type NewRoomRequest struct {
	Email       string `json:"email"`
	RoomName    string `json:"roomName"`
	MessageType string `bson:"-" json:"msgType"`
}

// File save files making sure they are distinct.
type File struct {
	MsgType        string `bson:"-" json:"msgType,omitempty"`
	UniqueFileHash string `bson:"_id" json:"fileHash"`
	FileName       string `bson:"fileName" json:"fileName"`
	User           string `bson:"userID" json:"userID"`
	FileSize       string `bson:"fileSize" json:"fileSize"`
	FileType       string `bson:"fileType" json:"fileType"`
	Chunks         int    `bson:"chunks,omitempty" json:"chunks"`
}

type FileChunks struct {
	MsgType            string `bson:"-" json:"msgType,omitempty"`
	FileName           string `bson:"-" json:"fileName,omitempty"`
	UniqueFileHash     string `bson:"_id" json:"fileHash,omitempty"`
	CompressedFileHash string `bson:"compressedFileHash" json:"compressedFileHash,omitempty"`
	FileBinary         string `bson:"fileChunk" json:"fileChunk,omitempty"`
	ChunkIndex         int    `bson:"chunkIndex" json:"chunkIndex"`
}

type WSMessage struct {
	Data []byte
	User string
}

// Connection is an middleman between the websocket connection and the hub.
type Connection struct {
	WS *websocket.Conn

	Send chan []byte
}

type Subscription struct {
	Conn *Connection
	User string
}

// Hub maintains the set of active connections and broadcasts messages to the
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

// classSessionPeerConnections allows class session video calls where a publisher
// broadcasts audio and video to other subscribers. Other subscribers broadcasts
// audio to each other.
// Video is mapped to unique ID which is to be generated by server for current class session
// to video and audio track.
// If UUID to current VIDEO track is nil. It indicates the room session is over.
// If publisher logs off, all peer connections related to that room is closed.
// Mutexes is integrated with video, audio and peerconnections to ensure data race free.
type classSessionPeerConnections struct {
	api *webrtc.API

	publisherVideoTracks  map[string]*webrtc.Track // mapped sessionID to track
	publisherTrackMutexes *sync.RWMutex

	audioTracks       map[string][]*webrtc.Track // mapped sessionID to track
	audioTrackSender  map[*webrtc.Track][]rtpSenderData
	audioTrackMutexes *sync.RWMutex

	peerConnection        map[string]*webrtc.PeerConnection // peerConnection is mapped user to peerconnection.
	peerConnectionMutexes *sync.RWMutex

	connectedUsers      map[string][]string // publisher is mapped sessionID to all connected users.
	connectedUsersMutex *sync.RWMutex
}

// rtpSenderData saves user RTPSender.
// On remove track, users can easily map all audio track to its senders.
// userID maps sender to it's peerconnection.
type rtpSenderData struct {
	userID string
	sender *webrtc.RTPSender
}

type webmWriter struct {
	fileName                       string
	audioWriter, videoWriter       webm.BlockWriteCloser
	audioBuilder, videoBuilder     *samplebuilder.SampleBuilder
	audioTimestamp, videoTimestamp uint32
}

type sdpConstruct struct {
	MsgType        string `json:"msgType"`
	ClassSessionID string `json:"sessionID"`
	Author         string `json:"author"`
	AuthorName     string `json:"name"`
	UserID         string `json:"userID"`
	RoomID         string `json:"roomID"`
	SDP            string `json:"sdp"`

	peerConnection *webrtc.PeerConnection
}

type dropboxUploader struct {
	uploadClient       files.Client
	sharableLinkClient sharing.Client
	fileUploadInfo     *files.CommitInfo

	file         *os.File
	fileFullPath string
	fileSize     int64
}
