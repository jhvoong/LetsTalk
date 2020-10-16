<template>
  <v-row no-gutters class="fill-height">
    <v-col cols="1">
      <SideBar
        :activateAddRoomDialog="activateAddRoomDialog"
        :activateNotificationDialog="activateNotificationDialog"
        :deactivateAllDialogs="deactivateAllDialogs"
        :unreadNotifications="unreadNotificationsCount"
        :unreadRoomMessages="unreadRoomMessageCount"
      />
    </v-col>

    <v-col cols="11">
      <v-row
        style="height: 100%"
        align="center"
        justify="center"
        v-if="showNotificationDialog || showAddRoomDialog"
      >
        <v-card
          v-if="showNotificationDialog"
          height="min-content"
          outlined
          shaped
        >
          <v-row class="mx-5">
            <v-col cols="12">
              <v-card-title>
                <b>Notifications</b>
              </v-card-title>
            </v-col>

            <v-col cols="12">
              <v-divider></v-divider>
            </v-col>

            <v-col cols="12">
              <v-card-text>
                <v-container
                  fluid
                  style="max-height: 72vh"
                  class="overflow-y-auto"
                >
                  <v-card-subtitle v-if="joinRequests.length == 0"
                    >No Notifications</v-card-subtitle
                  >
                  <v-card
                    v-for="(joinRequest, index) in joinRequests"
                    :key="index"
                    flat
                  >
                    <v-card-text
                      >{{ joinRequest.requestingUserName }} [{{
                        joinRequest.requestingUserID
                      }}] wants you to join room
                      {{ joinRequest.roomName }}</v-card-text
                    >
                    <v-card-actions>
                      <v-btn
                        @click="
                          joinRoom(
                            joinRequest.requestingUserID,
                            joinRequest.roomID,
                            joinRequest.roomName,
                            true,
                            index
                          )
                        "
                        text
                        color="green"
                        >Join</v-btn
                      >

                      <v-btn
                        @click="
                          joinRoom(
                            joinRequest.requestingUserID,
                            joinRequest.roomID,
                            joinRequest.roomName,
                            false,
                            index
                          )
                        "
                        text
                        color="red"
                        >Reject</v-btn
                      >
                    </v-card-actions>
                  </v-card>
                </v-container>
              </v-card-text>
            </v-col>
          </v-row>
        </v-card>

        <v-card flat outlined shaped v-else-if="showAddRoomDialog">
          <v-row class="mx-5">
            <v-col cols="12">
              <v-card-text>
                <b>Create New Room</b>
              </v-card-text>
            </v-col>

            <v-col cols="12">
              <v-text-field
                @keyup.enter.exact="createRoom()"
                placeholder="Room Name"
                rounded
                outlined
                v-model="newRoomName"
              />
              <v-btn @click="createRoom()" text color="green"
                >Create Room</v-btn
              >
            </v-col>
          </v-row>
        </v-card>
      </v-row>

      <v-row v-else no-gutters>
        <v-col cols="3">
          <RoomsPage
            :indexOfCurrentViewedRoom="indexOfCurrentViewedRoom"
            :recentChatPreview="recentChatPreview"
            :changeViewedRoomIndex="changeViewedRoomIndex"
            :joinedRooms="joinedRooms"
            :sendWSMessage="sendWSMessage"
            :unreadRoomMessages="unreadRoomMessages"
          />
        </v-col>

        <v-col cols="9">
          <ChatPage
            :initiateFile="initiateFile"
            :fileUploadDownload="fileUploadDownload"
            :associates="usersOnline"
            v-if="showChatPage"
            :clearFetchedUsers="clearFetchedUsers"
            :fetchedUsers="fetchedUsers"
            :sendWSMessage="sendWSMessage"
            :currentViewedRoom="currentViewedRoom"
            :changeDownloadStatus="changeDownloadStatus"
          />
        </v-col>
      </v-row>
    </v-col>
  </v-row>
</template>

<script lang="ts">
"use strict";

let socket: WebSocket;
// @ is an alias to /src
import Vue from "vue";
import store from "@/store";
import router from "@/router";
import CryptoJS from "crypto-js";

import { WSMessageType, MessageType, DefaultChunkSize } from "./Constants";

import SideBar from "../components/SideBar.vue";
import RoomsPage from "../components/RoomsPage.vue";
import ChatPage from "../components/ChatPage.vue";

import {
  JoinedRoom,
  RoomPageDetails,
  Message,
  JoinRequest,
  SentRoomRequest,
  FetchedUsers,
  UnreadRooms,
  RecentChatPreview,
  UsersOnline,
  ExitRoomDetails,
  FileUploadDownloadDetails,
  FileUploadComplete,
  FileDownload,
} from "./Types";

export default Vue.extend({
  name: "Home",
  components: {
    SideBar,
    RoomsPage,
    ChatPage,
  },

  data: () => ({
    joinedRooms: [] as JoinedRoom[],
    currentViewedRoom: {} as RoomPageDetails,
    joinRequests: [] as JoinRequest[],
    fetchedUsers: [] as FetchedUsers[],
    unreadRoomMessages: {} as UnreadRooms,
    recentChatPreview: {} as RecentChatPreview,
    usersOnline: {} as UsersOnline,
    fileUploadDownload: {} as FileDownload,

    userID: store.state.email,
    newRoomName: "",
    currentDownloadRoomID: "",

    showAddRoomDialog: false,
    showNotificationDialog: false,
    showChatPage: false,
    isFile: false,

    indexOfCurrentViewedRoom: 0,
    unreadRoomMessageCount: 0,
    unreadNotificationsCount: 0,
  }),

  methods: {
    onUnAuthorizedAccess: function () {
      router.push("/login");
      socket.close();
    },

    onWebsocketOpen: function (
      joinedRooms: JoinedRoom[],
      joinRequests: JoinRequest[]
    ) {
      if (joinedRooms) {
        this.joinedRooms = joinedRooms;
      }

      if (joinRequests) {
        this.joinRequests = joinRequests;
        this.unreadNotificationsCount = joinRequests.length;
        console.log(this.unreadNotificationsCount);
      }
    },

    onRequestMessages: function (roomDetails: RoomPageDetails) {
      if (roomDetails.firstLoad) {
        console.log(roomDetails.registeredUsers);
        this.currentViewedRoom = roomDetails;
        if (!this.currentViewedRoom.messages) {
          this.currentViewedRoom.messages = [];
        }

        this.joinedRooms[
          this.indexOfCurrentViewedRoom
        ].roomIcon = this.currentViewedRoom.roomIcon;

        this.scrollToBottomOfChatPage();
      } else {
        for (let i = roomDetails.messages.length - 1; i >= 0; i--) {
          this.currentViewedRoom.messages.unshift(roomDetails.messages[i]);
        }
      }
    },

    onJoinRoom: function (joinedRoom: JoinedRoom) {
      if (joinedRoom.userID === this.userID) {
        this.joinedRooms.unshift(joinedRoom);
        this.indexOfCurrentViewedRoom++;
        return;
      } else if (this.currentViewedRoom.roomID === joinedRoom.roomID) {
        let message: string = joinedRoom.userID + " Rejected join request.";
        if (joinedRoom.joined) {
          message = joinedRoom.userID + " Accepted join request.";
        }

        this.currentViewedRoom.messages.push({
          time: "",
          type: MessageType.Info,
          name: "",
          userID: "",
          roomID: "",
          size: "",
          hash: "",
          message: message,
          index: this.currentViewedRoom.messages.length + 1,
        });
      }

      this.usersOnline[joinedRoom.userID] = {
        name: joinedRoom.name,
        isOnline: true,
      };
    },

    onNewMessage: function (message: Message) {
      this.updateRecentMessagePreview(message.roomID, message.message);
      if (this.currentViewedRoom.roomID === message.roomID) {
        this.currentViewedRoom.messages.push(message);
        this.scrollToBottomOfChatPage();
        return;
      }

      this.changeNumberOfUnreadNotification(message.roomID, true);
      // ToDo: add notification sound
    },

    onSentRoomRequest: function (sentRequest: SentRoomRequest) {
      const message =
        sentRequest.userRequested +
        " Was requested to join the room by " +
        sentRequest.requesterID;

      if (this.userID === sentRequest.userRequested) {
        this.joinRequests.unshift({
          requestingUserName: sentRequest.requesterName,
          requestingUserID: sentRequest.requesterID,
          roomID: sentRequest.roomID,
          roomName: sentRequest.roomName,
        });

        this.unreadNotificationsCount = this.joinRequests.length;
        return;
      } else if (this.currentViewedRoom.roomID == sentRequest.roomID) {
        this.currentViewedRoom.messages.push({
          type: MessageType.Info,
          message: message,
          size: "",
          hash: "",
          time: "",
          name: "",
          userID: "",
          roomID: "",
          index: this.currentViewedRoom.messages.length + 1,
        });

        console.log("sent notification");
      }

      this.updateRecentMessagePreview(sentRequest.roomID, message);
    },

    onExitRoom: function (exitRoom: ExitRoomDetails) {
      console.log(exitRoom);
      const message: Message = {
        time: "",
        name: "",
        size: "",
        hash: "",
        type: MessageType.Info,
        userID: "",
        roomID: "",
        message: "You left the room.",
        index: this.currentViewedRoom.messages.length + 1,
      };

      if (exitRoom.roomID == this.currentViewedRoom.roomID) {
        if (exitRoom.userID != this.userID)
          message.message = exitRoom.userID + " Left the room.";

        this.currentViewedRoom.messages.push(message);
      }
    },

    onUploadFileChunks: function (file: FileUploadDownloadDetails) {
      if (
        this.fileUploadDownload[file.fileHash] &&
        !this.fileUploadDownload[file.fileHash].isDownloader
      ) {
        const fileDetails = this.fileUploadDownload[file.fileHash];
        fileDetails.chunk = file.nextChunk;

        // If file is successfully downloaded, send a FileUploadSuccess message to server
        // so as to be broadcasted to other users.
        if (fileDetails.chunk == fileDetails.chunks) {
          fileDetails.downloading = false;
          fileDetails.progress = 100;
          console.log("success");

          const message = {
            msgType: WSMessageType.FileUploadSuccess,
            userID: this.userID,
            name: this.userID,
            fileName: fileDetails.fileName,
            roomID: fileDetails.roomID,
            fileSize: fileDetails.fileSize.toString() + "MB",
            fileHash: fileDetails.fileHash,
          };
          socket.send(JSON.stringify(message));

          const infoMessage: Message = {
            time: "",
            size: "",
            hash: "",
            type: MessageType.Info,
            name: "",
            userID: "",
            roomID: fileDetails.roomID,
            message: "File successfully uploaded.",
            index: 0,
          };

          delete this.fileUploadDownload[fileDetails.fileHash];

          if (this.currentViewedRoom.roomID === fileDetails.roomID)
            this.currentViewedRoom.messages.push(infoMessage);

          return;
        }

        fileDetails.progress = (fileDetails.chunk / fileDetails.chunks) * 100;
        this.updateChatRoomPage();

        if (fileDetails.downloading === false) {
          return;
        }

        const offset = fileDetails.chunk * DefaultChunkSize;

        const reader = new FileReader();
        reader.onload = (event: ProgressEvent<FileReader>) => {
          if (!event.target) {
            return;
          }

          const data = event.target.result;
          if (!data) {
            return;
          }

          const fileUniqueHash = CryptoJS.SHA256(data.toString()).toString();

          const message = {
            msgType: WSMessageType.UploadFileChunk,
            userID: this.userID,
            file: data,
            fileName: fileDetails.fileName,
            fileHash: fileDetails.fileHash,
            newChunkHash: fileUniqueHash,
            recentChunkHash: "",
            chunkIndex: fileDetails.chunk,
          };

          console.log("sending", fileDetails.chunk, "now");

          // Send new file upload information to server, chunk == 0 indicates a new file.
          if (fileDetails.chunk === 0) {
            socket.send(JSON.stringify(message));
            return;
          }

          const recentOffset = (fileDetails.chunk - 1) * DefaultChunkSize;

          // Get recent file chunk hash.
          const recentFileChunkReader = new FileReader();
          recentFileChunkReader.onload = (e: ProgressEvent<FileReader>) => {
            if (!e.target) {
              return;
            }

            const recentData = e.target.result;
            if (!recentData) {
              return;
            }

            const recentFileUniqueHash = CryptoJS.SHA256(
              recentData.toString()
            ).toString();

            message.recentChunkHash = recentFileUniqueHash;
            socket.send(JSON.stringify(message));
          };

          if (fileDetails.fileContent) {
            const fileContent = fileDetails.fileContent;

            recentFileChunkReader.readAsDataURL(
              fileContent.slice(recentOffset, recentOffset + DefaultChunkSize)
            );
          }
        };

        if (fileDetails.fileContent) {
          const fileContent = fileDetails.fileContent;

          reader.readAsDataURL(
            fileContent.slice(offset, offset + DefaultChunkSize)
          );
        }
      }
    },

    onDownloadFileChunk: function (fileDetails: FileUploadDownloadDetails) {
      if (
        !fileDetails.compressedFileHash ||
        !this.fileUploadDownload[fileDetails.compressedFileHash]
      ) {
        console.log(
          "Server sent an invalid hash on file download",
          fileDetails
        );
        return;
      }

      if (!fileDetails.fileChunk) {
        console.log("Server sent a null file chunk.");
        return;
      }

      const file = this.fileUploadDownload[fileDetails.compressedFileHash];

      if (
        fileDetails.fileChunk &&
        fileDetails.fileHash !==
          CryptoJS.SHA256(fileDetails.fileChunk).toString()
      ) {
        console.log("File hash incorrect.");
        file.downloading = false;
        this.updateChatRoomPage();

        return;
      }

      fileDetails.fileChunk = fileDetails.fileChunk.substr(
        fileDetails.fileChunk.indexOf(",") + 1
      );

      const byteCharacters = atob(fileDetails.fileChunk);
      const byteArrays = [];

      for (
        let offset = 0;
        offset < byteCharacters.length;
        offset += DefaultChunkSize
      ) {
        const slice = byteCharacters.slice(offset, offset + DefaultChunkSize);
        const byteNumbers = new Array(slice.length);

        for (let i = 0; i < slice.length; i++) {
          byteNumbers[i] = slice.charCodeAt(i);
        }

        const byteArray = new Uint8Array(byteNumbers);
        byteArrays.push(byteArray);
      }

      const blob = new Blob(byteArrays, { type: "" });
      file.fileContent.push(blob);

      if (file.chunk == file.chunks) {
        this.downloadBlob(fileDetails.compressedFileHash);

        delete this.fileUploadDownload[fileDetails.compressedFileHash];
        this.updateChatRoomPage();
        return;
      }

      if (file.downloading === false) {
        console.log("file not downloading");
        return;
      }

      file.chunk++;
      file.progress = (file.chunk / file.chunks) * 100;

      this.updateChatRoomPage();
      const message = {
        userID: this.userID,
        msgType: WSMessageType.DownloadFileChunk,
        fileName: file.fileName,
        compressedFileHash: fileDetails.compressedFileHash,
        chunkIndex: file.chunk,
      };

      socket.send(JSON.stringify(message));
    },

    onRequestDownload: function (fileDetails: FileUploadDownloadDetails) {
      if (!fileDetails.fileHash) {
        console.log("File hash empty on request download", fileDetails);
        return;
      }

      const file = this.fileUploadDownload[fileDetails.fileHash];
      file.downloading = true;
      file.chunks = fileDetails.chunks;
      this.updateChatRoomPage();

      const message = {
        userID: this.userID,
        msgType: WSMessageType.DownloadFileChunk,
        fileName: file.fileName,
        // File hash of compressed file.
        compressedFileHash: file.fileHash,
        chunkIndex: file.chunk,
      };
      socket.send(JSON.stringify(message));
    },

    downloadBlob: function (fileHash: string) {
      const file = this.fileUploadDownload[fileHash];
      file.downloading = false;

      const blobB = new Blob(file.fileContent, {
        type: file.fileType,
      });

      const a = document.createElement("a");
      a.href = URL.createObjectURL(blobB);
      a.download = file.fileName;
      a.target = "_blank";
      a.click();
    },

    changeDownloadStatus: function (fileHash: string, isDownloading: boolean) {
      this.fileUploadDownload[fileHash].downloading = isDownloading;
      this.updateChatRoomPage();
    },

    broadcastFileToRoom: function (completeFileDetails: FileUploadComplete) {
      const message: Message = {
        time: "",
        size: completeFileDetails.fileSize,
        hash: completeFileDetails.fileHash,
        type: MessageType.File,
        name: completeFileDetails.name,
        userID: completeFileDetails.userID,
        roomID: completeFileDetails.roomID,
        message: completeFileDetails.fileName,
        index: this.currentViewedRoom.messages.length + 1,
      };

      if (this.currentViewedRoom.roomID === completeFileDetails.roomID) {
        this.currentViewedRoom.messages.push(message);
      }
    },

    initiateFile: function (file: FileUploadDownloadDetails) {
      this.fileUploadDownload[file.fileHash] = file;
      console.log(this.fileUploadDownload);
    },

    updateRecentMessagePreview: function (roomID: string, message: string) {
      this.recentChatPreview[roomID] = message;
    },

    updateChatRoomPage: function () {
      this.$nextTick(() => {
        setTimeout(() => {
          this.$children[2].$forceUpdate();
        }, 0);
      });
    },

    scrollToBottomOfChatPage: function () {
      this.$nextTick(() => {
        if (this.$children.length > 2) {
          const scrollHeight = this.$children[2].$el.querySelector("#messages");
          if (scrollHeight) scrollHeight.scrollTop = scrollHeight.scrollHeight;
        }
      });
    },

    changeNumberOfUnreadNotification: async function (
      roomID: string,
      addToUnreadNotifs: boolean
    ) {
      this.unreadRoomMessages[roomID] = addToUnreadNotifs;

      let messageCount = 0;
      for (const room in this.unreadRoomMessages) {
        if (
          this.unreadRoomMessages[room] &&
          this.unreadRoomMessages[room] == true
        )
          messageCount++;
      }

      this.unreadRoomMessageCount = messageCount;
    },

    changeViewedRoomIndex: function (index: number) {
      this.indexOfCurrentViewedRoom = index;
      this.changeNumberOfUnreadNotification(
        this.joinedRooms[index].roomID,
        false
      );
    },

    joinRoom: function (
      requestingUserID: string,
      roomID: string,
      roomName: string,
      joined: boolean,
      index: number
    ) {
      const message = {
        msgType: WSMessageType.JoinRoom,
        userID: this.userID,
        requesterID: requestingUserID,
        joined: joined,
        roomID: roomID,
        roomName: roomName,
      };

      socket.send(JSON.stringify(message));
      this.joinRequests.splice(index, 1);
      this.unreadNotificationsCount = this.joinRequests.length;
    },

    createRoom: function () {
      const message = {
        userID: this.userID,
        msgType: WSMessageType.CreateRoom,
        roomName: this.newRoomName,
      };

      socket.send(JSON.stringify(message));
      this.newRoomName = "";
    },

    activateAddRoomDialog: function () {
      this.showAddRoomDialog = true;
      this.showNotificationDialog = false;
    },

    clearFetchedUsers: function () {
      this.fetchedUsers = [];
    },

    activateNotificationDialog: function () {
      this.showNotificationDialog = true;
      this.showAddRoomDialog = false;
    },

    deactivateAllDialogs: function () {
      this.showAddRoomDialog = false;
      this.showNotificationDialog = false;
    },

    sendWSMessage: function (message: string) {
      socket.send(message);
    },
  },

  mounted() {
    // Verify login, if token is specified.
    if (store.state.token == "") {
      router.push("/login");
      return;
    }

    const URL: string = process.env.VUE_APP_BACKEND_SERVER;
    const URLs = URL.split("https://");

    // Probably write a test for this.
    if (URLs.length == 0) {
      router.push("/login");
      return;
    }

    socket = new WebSocket(
      "wss://" + URLs[1] + "/ws?token=" + store.state.token
    );

    socket.onmessage = (event: MessageEvent) => {
      const jsonContent = JSON.parse(event.data);

      switch (jsonContent.msgType) {
        case WSMessageType.UnauthorizedAccess:
          this.onUnAuthorizedAccess();
          break;

        case WSMessageType.SearchUser:
          this.fetchedUsers = jsonContent.fetchedUsers;
          console.log(this.fetchedUsers);
          break;

        case WSMessageType.WebsocketOpen:
          this.onWebsocketOpen(
            jsonContent.joinedRooms,
            jsonContent.joinRequests
          );
          this.usersOnline = jsonContent.associates;
          break;

        case WSMessageType.RequestMessages:
          if (!this.showChatPage) {
            this.showChatPage = true;
          }

          this.onRequestMessages(jsonContent.roomPageDetails);
          break;

        case WSMessageType.JoinRoom:
          this.onJoinRoom(jsonContent);
          break;

        case WSMessageType.NewMessage:
          this.onNewMessage(jsonContent);
          break;

        case WSMessageType.SentRoomRequest:
          this.onSentRoomRequest(jsonContent);
          break;

        case WSMessageType.OnlineStatus:
          if (this.usersOnline[jsonContent.userID]) {
            this.usersOnline[jsonContent.userID].isOnline = jsonContent.status;
          }
          break;

        case WSMessageType.ExitRoom:
          this.onExitRoom(jsonContent);
          break;

        case WSMessageType.UploadFileChunk:
          this.onUploadFileChunks(jsonContent);
          break;

        case WSMessageType.FileUploadSuccess:
          this.broadcastFileToRoom(jsonContent);
          break;

        case WSMessageType.UploadFileError:
          if (this.$children.length > 2) {
            this.$children[2].$emit("onUploadError");
          }
          break;

        case WSMessageType.FileRequestDownload:
          this.onRequestDownload(jsonContent);
          break;

        case WSMessageType.DownloadFileChunk:
          console.log("Downloading");
          this.onDownloadFileChunk(jsonContent);
          break;
      }
    };

    socket.onerror = (event: Event) => {
      // Reconnect Websocket if not UnAuthorized.
      console.log("Websocket errored.", event);
    };

    socket.onopen = () => {
      console.log("Websocket open.");

      // Fetch users content from API.
      // Contents that are to be fetched from API are, Registered rooms and room request.
      const message = {
        msgType: WSMessageType.WebsocketOpen,
        userID: this.userID,
      };

      socket.send(JSON.stringify(message));
    };

    socket.onclose = function (closeEvent: CloseEvent) {
      // ToDo: We can create a dialog and ask users if they want to reconnect.
      console.log("Websocket closed.", closeEvent);
    };
  },
});
</script>


<style  scoped>
.transparent_sheet {
  opacity: 0.2;
  border-color: transparent !important;
  z-index: -1;
  position: absolute;
}
</style>