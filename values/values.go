package values

import (
	"github.com/pion/webrtc/v2"
	log "github.com/sirupsen/logrus"
)

type MessageType int

const (
	// Files are saved as base64 format to database then if queried,
	TXT MessageType = iota
	INFO
	FILE
)

// All request message types both clients and server
const (
	NewFileUploadMsgType          = "NewFileUpload"
	NewMessageMsgType             = "NewMessage"
	RequestAllMessagesMsgType     = "RequestAllMessages"
	SearchUserMsgType             = "SearchUser"
	WebsocketOpenMsgType          = "WebsocketOpen"
	NewRoomCreatedMsgType         = "NewRoomCreated"
	ExitRoomMsgType               = "ExitRoom"
	RequestUsersToJoinRoomMsgType = "RequestUsersToJoinRoom"
	UserJoinedRoomMsgType         = "UserJoinedRoom"

	UploadFileErrorMsgType   = "UploadFileError" // UploadFileErrorMsgType is sent to client only.
	UploadFileSuccessMsgType = "FileUploadSuccess"
	UploadFileChunkMsgType   = "UploadFileChunk"

	RequestDownloadMsgType     = "RequestDownload"
	DownloadFileChunkMsgType   = "DownloadFileChunk"
	DownloadFileErrorMsgType   = "DownloadFileError"   // DownloadFileErrorMsgType is sent to client only.
	DownloadFileSuccessMsgType = "DownloadFileSuccess" // DownloadFileSuccessMsgType is sent to client only.

	StartClassSession = "StartClassSession"
	JoinClassSession  = "JoinClassSession"
	NegotiateSDP      = "RenegotiateSDP"
	ClassSessionError = "ClassSessionError"
)

var (
	// MapEmailToName maps user email to name
	MapEmailToName map[string]string

	// PeerConnectionConfig contains peerconnection configuration
	PeerConnectionConfig = webrtc.Configuration{
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
	}

	credetialType map[string]webrtc.ICECredentialType = map[string]webrtc.ICECredentialType{
		"Password": webrtc.ICECredentialTypePassword,
		"Oauth":    webrtc.ICECredentialTypePassword,
	}
)

func initIceServers() {
	if len(Config.ICEServers) == 0 {
		PeerConnectionConfig.ICEServers = []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		}

		return
	}

	for _, config := range Config.ICEServers {
		if len(config.URLs) == 0 {
			log.Fatalln("User did not specify ICE server.")
		}

		credential, ok := credetialType[config.AuthType]
		if !ok {
			log.Fatalln("Invalid webrtc credential type", config.AuthType, "only AuthType Password and Oauth are allowed.")
		}

		iceServer := webrtc.ICEServer{
			URLs:           config.URLs,
			Username:       config.Username,
			Credential:     config.AuthSecret,
			CredentialType: credential,
		}

		PeerConnectionConfig.ICEServers = append(PeerConnectionConfig.ICEServers, iceServer)
	}
}
