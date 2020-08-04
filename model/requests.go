package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/metaclips/LetsTalk/values"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type messageBytes []byte

// handleCreateNewRoom creates a new room for user.
func (msg messageBytes) handleCreateNewRoom() {
	var newRoom NewRoomRequest
	if err := json.Unmarshal(msg, &newRoom); err != nil {
		log.Warningln("Could not convert to required New Room Request struct")
		return
	}

	roomID, err := newRoom.createNewRoom()
	if err != nil {
		log.Warningln("Unable to create a new room for user:", newRoom.Email, "err:", err.Error())
		return
	}

	// Broadcast a joined message.
	userJoinedMessage := Joined{
		RoomID:      roomID,
		Email:       newRoom.Email,
		RoomName:    newRoom.RoomName,
		MessageType: "UserJoinedRoom",
	}

	jsonContent, err := json.Marshal(userJoinedMessage)
	if err != nil {
		log.Warningln("Could not marshal to jsonByte while creating room", err.Error())
		return
	}

	HubConstruct.sendMessage(jsonContent, newRoom.Email)
}

func (msg messageBytes) handleRequestUserToJoinRoom() {
	var request JoinRequest
	if err := json.Unmarshal(msg, &request); err != nil {
		log.Warningln("Could not convert to required Joined Request struct")
		return
	}

	for _, user := range request.Users {
		roomRegisteredUser, err := request.requestUserToJoinRoom(user)
		if err != nil {
			log.Warningln("Error while requesting to room", err)
			continue
		}

		data := struct {
			RequesterID   string `json:"requesterID"`
			RequesterName string `json:"requesterName"`
			UserRequested string `json:"userRequested"`
			RoomID        string `json:"roomID"`
			RoomName      string `json:"roomName"`
			MsgType       string `json:"msgType"`
		}{
			request.RequestingUserID, request.RequestingUserName,
			user, request.RoomID, request.RoomName, "RequestUsersToJoinRoom",
		}

		jsonContent, err := json.Marshal(data)
		if err != nil {
			log.Warningln("could not marshal to RequestUsersToJoinRoom, err:", err)
			continue
		}

		// Send back RequestUsersToJoinRoom signal to everyone registered in room.
		for _, roomRegisteredUser := range roomRegisteredUser {
			HubConstruct.sendMessage(jsonContent, roomRegisteredUser)
		}

		HubConstruct.sendMessage(jsonContent, user)
	}
}

// handleUserAcceptRoomRequest accepts room join request.
func (msg messageBytes) handleUserAcceptRoomRequest(joiner string) {
	var roomRequest Joined
	if err := json.Unmarshal(msg, &roomRequest); err != nil {
		log.Warningln("Could not convert to required Join Room Request struct")
		return
	}

	if roomRequest.Email != joiner {
		return
	}

	users, err := roomRequest.acceptRoomRequest()
	if err != nil {
		log.Warningln("could not join room", err)
		return
	}

	for _, user := range users {
		HubConstruct.sendMessage(msg, user)
	}
}

// handleNewMessage broadcasts users message to all online users and also saves to database.
func (msg messageBytes) handleNewMessage(author string) {
	var newMessage Message
	if err := json.Unmarshal(msg, &newMessage); err != nil {
		log.Warningln("Could not convert to required New Message struct", err)
		return
	}

	newMessage.Time = time.Now().Format(values.TimeLayout)
	// Save message to database ensuring user is registered to room.
	registeredUsers, err := newMessage.saveMessageContent()
	if err != nil {
		log.Warningln("Error saving msg to db", err, author)
		return
	}

	jsonContent, err := json.Marshal(newMessage)
	if err != nil {
		log.Warningln("Error converted message to json content", err)
		return
	}

	// Message is sent back to all online users including sender.
	for _, registeredUser := range registeredUsers {
		HubConstruct.sendMessage(jsonContent, registeredUser)
	}
}

// handleExitRoom exits requesters joined room and also notifies all room users.
func (msg messageBytes) handleExitRoom(author string) {
	data := struct {
		Email  string `json:"email"`
		RoomID string `json:"roomID"`
	}{}

	if err := json.Unmarshal(msg, &data); err != nil {
		log.Warningln("Could not retrieve json on exit room request", err)
		return
	}

	if author != data.Email {
		return
	}

	user := User{Email: data.Email}
	registeredUsers, err := user.exitRoom(data.RoomID)
	if err != nil {
		log.Warningln("Error exiting room", err)
		return
	}

	// Broadcast to all online users of a room exit.
	for _, registeredUser := range registeredUsers {
		HubConstruct.sendMessage(msg, registeredUser)
	}

	HubConstruct.sendMessage(msg, author)
}

// handleNewFileUpload creates a new file content in database.
// If file create is a success, a file upload success is sent to client to send next chunk.
// next chunk could be the next preceding file chunk if another user has uploaded file content.
// If file upload error, send back error message to user
func (msg messageBytes) handleNewFileUpload() {
	file := File{}
	if err := json.Unmarshal(msg, &file); err != nil {
		log.Warningln(err)
		return
	}

	data := struct {
		MsgType      string `json:"msgType"`
		ErrorMessage string `json:"errorMsg,omitempty"`
		RecentHash   string `json:"recentHash"`
		FileName     string `json:"fileName,omitempty"`
		Chunk        int    `json:"nextChunk"`
	}{}

	data.FileName = file.FileName
	user := file.User

	if err := file.uploadNewFile(); err == mongo.ErrNoDocuments || err == nil {
		// Send next file chunk and current hash which is a "".
		data.MsgType = values.UploadFileChunkMsgType

		// Resume file chunk upload if Current chunk is greater than 0.
		if file.Chunks > 0 {
			data.Chunk = file.Chunks + 1
		} else {
			data.Chunk = file.Chunks
		}

	} else {
		log.Warningln("Error on handle new file upload calling UploadNewFile", err)
		data.ErrorMessage = values.ErrFileUpload.Error()
		data.MsgType = values.UploadFileErrorMsgType
	}

	jsonContent, err := json.Marshal(&data)
	if err != nil {
		log.Warningln("Error sending marshalled ")
		return
	}

	HubConstruct.sendMessage(jsonContent, user)
}

func (msg messageBytes) handleUploadFileChunk() {
	data := struct {
		MsgType            string `json:"msgType"`
		User               string `json:"userID"`
		FileName           string `json:"fileName"`
		File               string `json:"file,omitempty"`
		NewChunkHash       string `json:"newChunkHash,omitempty"`
		RecentChunkHash    string `json:"recentChunkHash,omitempty"`
		ChunkIndex         int    `json:"chunkIndex,omitempty"`
		NextChunk          int    `json:"nextChunk"`
		CompressedFileHash string `json:"compressedFileHash"`
	}{}

	if err := json.Unmarshal(msg, &data); err != nil {
		fmt.Println(string(msg))
		log.Warningln("error marshalling handle file upload struct", err)
		return
	}

	file := FileChunks{
		UniqueFileHash:     data.NewChunkHash,
		FileBinary:         data.File,
		ChunkIndex:         data.ChunkIndex,
		CompressedFileHash: data.CompressedFileHash,
	}

	userID := data.User
	var recentFileExist bool
	// If file upload is a new file, set recent file exist as true.
	if data.RecentChunkHash == "" {
		recentFileExist = true
	} else {
		recentFileExist = FileChunks{UniqueFileHash: data.RecentChunkHash}.fileChunkExists()
	}

	data.RecentChunkHash, data.File, data.NewChunkHash = "", "", ""
	data.NextChunk, data.ChunkIndex = 0, 0

	fileHash := sha256.Sum256([]byte(file.FileBinary))
	// Check if client sent file hash is same as server generated Hash.
	if hex.EncodeToString(fileHash[:]) != file.UniqueFileHash || !recentFileExist {
		data.MsgType = "UploadError"

		// Re-request for current chunk index.
		jsonContent, err := json.Marshal(&data)
		if err != nil {
			log.Warningln("Could not generate jsonContent to re-request file chunk")
			return
		}

		HubConstruct.sendMessage(jsonContent, data.User)

		return
	}

	if err := file.addFileChunk(); err != nil {
		// What could be cases where err is not nil.
		// File could have already been added to database?.
		// We still request for next file chunk, if when we receive a new fille chunk,
		// so that when we notice file corruption, we re-request from corrupted stage.
		log.Warningln("error adding file chunk to server", err)
	}

	data.NextChunk = file.ChunkIndex + 1

	jsonContent, err := json.Marshal(&data)
	if err != nil {
		log.Warningln("error sending marshalled data to server on handle file upload chunk")
		return
	}

	HubConstruct.sendMessage(jsonContent, userID)
}

// handleUploadFileUploadComplete is called when file chunk uploads is complete.
// File accessibility is broadcasted to other users in the room so as to download file.
func (msg messageBytes) handleUploadFileUploadComplete() {
	data := struct {
		MsgType  string `json:"msgType"`
		UserID   string `json:"userID"`
		UserName string `json:"name"`
		FileName string `json:"fileName"`
		FileSize string `json:"fileSize"`
		FileHash string `json:"fileHash"`
		RoomID   string `json:"roomID"`
	}{}

	if err := json.Unmarshal(msg, &data); err != nil {
		log.Warningln("error unmarshalling struct for handleUploadFileUploadComplete", err)
		return
	}

	data.MsgType = values.UploadFileSuccessMsgType

	roomUsers, err := Message{
		RoomID:   data.RoomID,
		UserID:   data.UserID,
		Name:     data.UserName,
		Message:  data.FileName,
		Time:     time.Now().Format(values.TimeLayout),
		Type:     "file",
		FileSize: data.FileSize,
		FileHash: data.FileHash,
	}.saveMessageContent()

	if err != nil {
		log.Warningln("error saving message content on handleUploadFileUploadComplete", err)
	}

	jsonContent, err := json.Marshal(&data)
	if err != nil {
		log.Warningln("error marshalling json in handleUploadFileUploadComplete", err)
	}

	for _, roomUser := range roomUsers {
		if roomUser == data.UserID {
			continue
		}

		HubConstruct.sendMessage(jsonContent, roomUser)
	}
}

func (msg messageBytes) handleRequestDownload(author string) {
	file := File{}
	if err := json.Unmarshal(msg, &file); err != nil {
		log.Warningln("error unmarshalling json on handleRequestDownload", err)
		return
	}

	fileName := file.FileName

	if err := file.retrieveFileInformation(); err != nil {
		file.MsgType = values.DownloadFileErrorMsgType
	}
	file.FileName = fileName

	jsonContent, err := json.Marshal(&file)
	if err != nil {
		log.Warningln("error marshalling json on handleRequestDownload", err)
	}

	HubConstruct.sendMessage(jsonContent, author)
}

func (msg messageBytes) handleFileDownload(author string) {
	file := FileChunks{}
	if err := json.Unmarshal(msg, &file); err != nil {
		log.Warningln("error unmarshalling json on handleFileDownload", err)
		return
	}

	fileName := file.FileName

	if err := file.retrieveFileChunk(); err != nil {
		log.Warningln("error retrieving file", err)
		// Send download file error message to client so as to stop download.
		file = FileChunks{}
		file.MsgType = values.DownloadFileErrorMsgType
	} else {
		file.MsgType = values.DownloadFileChunkMsgType
	}

	file.FileName = fileName

	jsonContent, err := json.Marshal(&file)
	if err != nil {
		log.Warningln("error marshalling json on handleFileDownload", err)
	}

	HubConstruct.sendMessage(jsonContent, author)
}

// handleSearchUser returns registered users that match searchText.
func handleSearchUser(searchText, user string) {
	data := struct {
		UsersFound []string
		MsgType    string `json:"msgType"`
	}{
		GetUser(searchText, user),
		"getUsers",
	}

	jsonContent, err := json.Marshal(&data)
	if err != nil {
		log.Warningln("error while converting search user result to json", err)
		return
	}

	HubConstruct.sendMessage(jsonContent, user)
}

// handleRequestAllMessages coallates all messages in a particular room.
func handleRequestAllMessages(roomID, author string) {
	room := Room{RoomID: roomID}
	if err := room.getAllMessageInRoom(); err != nil {
		log.Warningln("could not get all messages in room, err:", err)
		return
	}

	onlineUsers := make(map[string]bool)
	for _, user := range room.RegisteredUsers {
		if name, ok := values.MapEmailToName[user]; ok {
			nameAndEmail := fmt.Sprintf("%s (%s)", name, user)
			onlineUsers[nameAndEmail] = HubConstruct.Users[user] != nil
		}
	}

	data := map[string]interface{}{
		"messages":    room.Messages,
		"msgType":     "RequestAllMessages",
		"roomName":    room.RoomName,
		"roomID":      roomID,
		"onlineUsers": onlineUsers,
	}

	jsonContent, err := json.Marshal(&data)
	if err != nil {
		log.Warningln("could not marshal images, err:", err)
		return
	}

	HubConstruct.sendMessage(jsonContent, author)
}

// handleLoadUserContent loads all users contents on page load.
// All rooms joined and users requests are loaded through WS.
func handleLoadUserContent(email string) {
	userInfo := User{Email: email}
	if err := userInfo.getUser(); err != nil {
		log.Warningln("could not fetch users room", email)
		return
	}

	request := map[string]interface{}{
		"msgType":  "WebsocketOpen",
		"rooms":    userInfo.RoomsJoined,
		"requests": userInfo.JoinRequest,
	}

	if data, err := json.Marshal(request); err == nil && HubConstruct.Users[email] != nil {
		HubConstruct.sendMessage(data, email)
	}
}

// broadcastOnlineStatusToAllUserRoom broadcasts users availability status to all users joined rooms.
// Status are broadcasted timely.
// Since we are calling broadcastOnlineStatusToAllUserRoom from HubRun, we should call it in a goroutine so as
// not to block the hub channel
func broadcastOnlineStatusToAllUserRoom(userEmail string, online bool) {
	user := User{Email: userEmail}
	associates, err := user.getAllUsersAssociates()
	if err != nil {
		log.Warningln("could not get users associate", err)
		return
	}

	for _, assassociateEmail := range associates {
		nameAndEmail := fmt.Sprintf("%s (%s)", values.MapEmailToName[assassociateEmail], userEmail)
		msg := map[string]interface{}{
			"msgType":  "OnlineStatus",
			"username": nameAndEmail,
			"status":   online,
		}

		if data, err := json.Marshal(msg); err == nil {
			HubConstruct.sendMessage(data, assassociateEmail)
		}
	}
}
