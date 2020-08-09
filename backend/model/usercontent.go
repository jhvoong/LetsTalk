package model

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/metaclips/LetsTalk/backend/values"
)

func (b *User) getUser() error {
	result := db.Collection(values.UsersCollectionName).FindOne(ctx, bson.M{
		"_id": b.Email,
	})

	if err := result.Decode(&b); err != nil {
		return err
	}

	return nil
}

func (b User) addUserToRoom(roomID, roomName string) error {
	b.updateRoomsJoinedByUsers(roomID, roomName)
	var chats Room
	message := Message{
		Message: b.Email + " Joined",
		Type:    getContentType(values.INFO),
	}

	result := db.Collection(values.RoomsCollectionName).FindOne(ctx, bson.M{
		"_id": roomID,
	})

	if err := result.Decode(&chats); err != nil {
		return err
	}

	chats.Messages = append(chats.Messages, message)
	_, err := db.Collection(values.RoomsCollectionName).UpdateOne(ctx, bson.M{"_id": roomID},
		bson.M{"$set": bson.M{"messages": chats.Messages}})

	return err
}

func (b *User) updateRoomsJoinedByUsers(roomID, roomName string) error {
	if err := b.getUser(); err != nil {
		return err
	}

	var roomJoined = RoomsJoined{RoomID: roomID, RoomName: roomName}
	b.RoomsJoined = append(b.RoomsJoined, roomJoined)

	_, err := db.Collection(values.UsersCollectionName).UpdateOne(ctx, bson.M{"_id": b.Email},
		bson.M{"$set": bson.M{"roomsJoined": b.RoomsJoined}})

	return err
}

func (b *User) getAllUsersAssociates() ([]string, error) {
	if err := b.getUser(); err != nil {
		return nil, err
	}

	usersChannel := make(chan []string)
	done := make(chan struct{})
	users := make([]string, 0)
	registeredUser := make(map[string]bool)

	go func() {
		for {
			data, ok := <-usersChannel
			if ok {
				for _, user := range data {
					if _, exist := registeredUser[user]; !exist && user != b.Email {
						users = append(users, user)
						registeredUser[user] = true
					}
				}

				continue
			}

			close(done)
			break
		}
	}()

	for _, roomJoined := range b.RoomsJoined {
		var room Room
		result := db.Collection(values.RoomsCollectionName).FindOne(ctx, bson.M{
			"_id": roomJoined.RoomID,
		})

		if err := result.Decode(&room); err != nil {
			close(usersChannel)
			<-done
			return nil, err
		}

		usersChannel <- room.RegisteredUsers
	}

	close(usersChannel)
	<-done

	return users, nil
}

func (b User) exitRoom(roomID string) ([]string, error) {
	if err := b.getUser(); err != nil {
		return nil, err
	}

	// Confirm if indeed user is registered to room
	var roomExist bool
	for i, roomJoined := range b.RoomsJoined {
		if roomJoined.RoomID == roomID {
			roomExist = true

			if len(b.RoomsJoined)-1 > i {
				b.RoomsJoined = append(b.RoomsJoined[:i], b.RoomsJoined[i+1:]...)
			} else {
				b.RoomsJoined = b.RoomsJoined[:i]
			}
			break
		}
	}
	if !roomExist {
		return nil, values.ErrUserNotRegisteredToRoom
	}

	// Update room joined by user in DB.
	_, err := db.Collection(values.UsersCollectionName).UpdateOne(ctx, bson.M{"_id": b.Email},
		bson.M{"$set": bson.M{"roomsJoined": b.RoomsJoined}})
	if err != nil {
		return nil, err
	}

	room := Room{RoomID: roomID}
	result := db.Collection(values.RoomsCollectionName).FindOne(ctx, bson.M{"_id": room.RoomID})

	if err := result.Decode(&room); err != nil {
		return nil, err
	}

	exitMessage := Message{
		Type:    getContentType(values.INFO),
		Message: b.Email + " left the room",
	}
	room.Messages = append(room.Messages, exitMessage)

	for i, user := range room.RegisteredUsers {
		if user == b.Email {
			if len(room.RegisteredUsers)-1 > i {
				room.RegisteredUsers = append(room.RegisteredUsers[:i], room.RegisteredUsers[i+1:]...)
			} else {
				room.RegisteredUsers = room.RegisteredUsers[:i]
			}

			break
		}
	}

	_, err = db.Collection(values.RoomsCollectionName).UpdateOne(ctx, bson.M{
		"_id": room.RoomID,
	}, bson.M{"$set": bson.M{"registeredUsers": room.RegisteredUsers, "messages": room.Messages}})

	return room.RegisteredUsers, err
}

func (b User) CreateUserLogin(password string, w http.ResponseWriter) error {
	if err := b.getUser(); err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(b.Password, []byte(password)); err != nil {
		return err
	}

	err := CookieDetail{
		Email:      b.Email,
		Collection: values.UsersCollectionName,
		CookieName: values.UserCookieName,
		Path:       "/",
		Data: CookieData{
			Email: b.Email,
		},
	}.CreateCookie(w)

	return err
}

func (b User) validateUser(uniqueID string) error {
	if err := b.getUser(); err != nil {
		return err
	}

	if b.UUID != uniqueID {
		return values.ErrIncorrectUUID
	}
	return nil
}

func (b Message) saveMessageContent() ([]string, error) {
	var messages Room
	result := db.Collection(values.RoomsCollectionName).FindOne(ctx, bson.M{
		"_id": b.RoomID,
	})

	if err := result.Decode(&messages); err != nil {
		return nil, err
	}

	var userExists bool
	for _, user := range messages.RegisteredUsers {
		if b.UserID == user {
			userExists = true
			break
		}
	}

	if !userExists {
		return nil, values.ErrInvalidUser
	}

	messages.Messages = append(messages.Messages, b)
	_, err := db.Collection(values.RoomsCollectionName).UpdateOne(ctx, bson.M{"_id": b.RoomID},
		bson.M{"$set": bson.M{"messages": messages.Messages}})

	return messages.RegisteredUsers, err
}

func (b NewRoomRequest) createNewRoom() (string, error) {
	var chats Room
	message := Message{
		Message: b.Email + " Joined",
		Type:    getContentType(values.INFO),
	}

	chats.Messages = append(chats.Messages, message)
	chats.RoomID = uuid.New().String()
	chats.RoomName = b.RoomName
	chats.RegisteredUsers = append(chats.RegisteredUsers, b.Email)

	if _, err := db.Collection(values.RoomsCollectionName).InsertOne(ctx, chats); err != nil {
		return "", err
	}

	user := User{Email: b.Email}
	if err := user.updateRoomsJoinedByUsers(chats.RoomID, chats.RoomName); err != nil {
		return "", err
	}

	return chats.RoomID, nil
}

func (b *Room) getAllMessageInRoom() error {
	result := db.Collection(values.RoomsCollectionName).FindOne(ctx, bson.M{"_id": b.RoomID})

	if err := result.Decode(&b); err != nil {
		return err
	}

	return nil
}

func (b Joined) acceptRoomRequest() ([]string, error) {
	result := db.Collection(values.UsersCollectionName).FindOne(ctx, bson.M{
		"_id": b.Email,
	})

	var user User
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}

	// Check users join requests for room.
	var joinRequestLegit bool
	for i, request := range user.JoinRequest {
		if request.RoomID == b.RoomID {
			joinRequestLegit = true
			user.JoinRequest = append(user.JoinRequest[:i], user.JoinRequest[i+1:]...)
			break
		}
	}

	if !joinRequestLegit {
		return nil, values.ErrIllicitJoinRequest
	}

	user.RoomsJoined = append(user.RoomsJoined, RoomsJoined{RoomID: b.RoomID, RoomName: b.RoomName})

	_, err = db.Collection(values.UsersCollectionName).UpdateOne(ctx, bson.M{"_id": b.Email},
		bson.M{"$set": bson.M{"joinRequest": user.JoinRequest, "roomsJoined": user.RoomsJoined}})
	if err != nil {
		return nil, err
	}

	result = db.Collection(values.RoomsCollectionName).FindOne(ctx, bson.M{
		"_id": b.RoomID,
	})

	var messages Room
	if err := result.Decode(&messages); err != nil {
		return nil, err
	}

	message := Message{
		Message: b.Email + " Joined",
		Type:    getContentType(values.INFO),
	}

	messages.RegisteredUsers = append(messages.RegisteredUsers, b.Email)
	messages.Messages = append(messages.Messages, message)

	_, err = db.Collection(values.RoomsCollectionName).UpdateOne(ctx, bson.M{
		"_id": b.RoomID,
	}, bson.M{"$set": bson.M{"registeredUsers": messages.RegisteredUsers, "messages": messages.Messages}})

	return messages.RegisteredUsers, err
}

func (b JoinRequest) requestUserToJoinRoom(userToJoinEmail string) ([]string, error) {
	var room Room
	result := db.Collection(values.RoomsCollectionName).FindOne(ctx, bson.M{"_id": b.RoomID})

	if err := result.Decode(&room); err != nil {
		return nil, err
	}

	// Confirm if person making the request is part of the room.
	var requesterLegit bool
	for _, registeredUser := range room.RegisteredUsers {
		if registeredUser == b.RequestingUserID {
			requesterLegit = true
			break
		} else if registeredUser == userToJoinEmail {
			return nil, values.ErrUserExistInRoom
		}
	}

	if !requesterLegit {
		return nil, errors.New("Invalid user made a RequestUsersToJoinRoom request Name: " + b.RequestingUserID)
	}

	result = db.Collection(values.UsersCollectionName).FindOne(ctx, bson.M{"_id": userToJoinEmail})
	var user User

	if err := result.Decode(&user); err != nil {
		return nil, err
	}

	// Check if user has already been requested by the room.
	for _, request := range user.JoinRequest {
		if b.RoomID == request.RoomID {
			return nil, values.ErrUserAlreadyRequested
		}
	}
	user.JoinRequest = append(user.JoinRequest, b)

	_, err := db.Collection(values.UsersCollectionName).UpdateOne(ctx, bson.M{"_id": userToJoinEmail},
		bson.M{"$set": bson.M{"joinRequest": user.JoinRequest}})

	if err != nil {
		return nil, err
	}

	message := Message{
		Message: fmt.Sprintf("%s was requested to join the room by %s", userToJoinEmail, b.RequestingUserID),
		Type:    getContentType(values.INFO),
	}
	room.Messages = append(room.Messages, message)

	_, err = db.Collection(values.RoomsCollectionName).
		UpdateOne(ctx, bson.M{"_id": b.RoomID}, bson.M{"$set": bson.M{"messages": room.Messages}})
	return room.RegisteredUsers, err
}

// UploadNewFile create a NewFile content to database and returns file content if one
// has already been created.
// Chunks is set to zero so that if user wants to retrieve
func (b *File) uploadNewFile() error {
	result := db.Collection(values.FilesCollectionName).FindOne(ctx, bson.M{"_id": b.UniqueFileHash}) //, b, options.FindOneAndReplace().SetUpsert(true))

	if result.Err() == mongo.ErrNoDocuments {
		_, err := db.Collection(values.FilesCollectionName).InsertOne(ctx, b)
		return err
	}

	if err := result.Decode(&b); err != nil {
		return err
	}

	return nil
}

func (b *File) retrieveFileInformation() error {
	result := db.Collection(values.FilesCollectionName).FindOne(ctx, bson.M{"_id": b.UniqueFileHash})
	return result.Decode(&b)
}

func (b FileChunks) fileChunkExists() bool {
	result := db.Collection(values.FileChunksCollectionName).FindOne(ctx, bson.M{"_id": b.UniqueFileHash})
	if err := result.Err(); err == nil {
		return true
	}
	return false
}

func (b FileChunks) addFileChunk() error {
	result := db.Collection(values.FileChunksCollectionName).
		FindOneAndReplace(ctx, bson.M{"_id": b.UniqueFileHash}, b, options.FindOneAndReplace().SetUpsert(true))

	// Update original file index.
	if err := result.Err(); err == nil || err == mongo.ErrNoDocuments {
		_, err := db.Collection(values.FilesCollectionName).UpdateOne(ctx,
			bson.M{"_id": b.CompressedFileHash}, bson.M{"$set": bson.M{"chunks": b.ChunkIndex}})
		return err
	}

	return result.Err()
}

func (b *FileChunks) retrieveFileChunk() error {
	result := db.Collection(values.FileChunksCollectionName).
		FindOne(ctx, bson.M{"compressedFileHash": b.CompressedFileHash, "chunkIndex": b.ChunkIndex})

	return result.Decode(&b)
}

func uploadFileGridFS(fileName string) error {
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println("unable read file while uploading", err)
		return err
	}

	buc, err := gridfs.NewBucket(db)
	if err != nil {
		log.Println("unable GridFS bucket", err)
		return err
	}

	up, err := buc.OpenUploadStream("hhh")
	if err != nil {
		log.Println("unable to open upload stream", err)
		return err
	}
	defer up.Close()

	_, err = up.Write(fileBytes)
	if err != nil {
		log.Println("unable to write to bucket stream", err)
		return err
	}

	return nil
}

func getContentType(contentType values.MessageType) string {
	switch contentType {
	case values.INFO:
		return "info"
	case values.TXT:
		return "txt"
	case values.FILE:
		return "file"
	}

	return ""
}

func GetUser(key string, user string) (names []string) {
	names = make([]string, 0)
	for email := range values.MapEmailToName {
		if email == "" || email == user {
			continue
		}
		if strings.Contains(email, key) {
			names = append(names, email)
		}
	}

	return
}
