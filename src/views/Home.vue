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
        style="height: 100%;"
        align="center"
        justify="center"
        v-if="showNotificationDialog||showAddRoomDialog"
      >
        <v-card v-if="showNotificationDialog" height="min-content" outlined shaped>
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
                <v-container fluid style="max-height: 72vh;" class="overflow-y-auto">
                  <v-card-subtitle v-if="joinRequests.length==0">No Notifications</v-card-subtitle>
                  <v-card v-for="(joinRequest, index) in joinRequests" :key="index" flat>
                    <v-card-text>{{joinRequest.requestingUserName}} [{{joinRequest.requestingUserID}}] wants you to join room {{joinRequest.roomName}}</v-card-text>
                    <v-card-actions>
                      <v-btn
                        @click="joinRoom(joinRequest.requestingUserID, joinRequest.roomID, joinRequest.roomName,true, index)"
                        text
                        color="green"
                      >Join</v-btn>

                      <v-btn
                        @click="joinRoom(joinRequest.requestingUserID, joinRequest.roomID, joinRequest.roomName,false, index)"
                        text
                        color="red"
                      >Reject</v-btn>
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
              <v-btn @click="createRoom()" text color="green">Create Room</v-btn>
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
            :associates="usersOnline"
            v-if="showChatPage"
            :clearFetchedUsers="clearFetchedUsers"
            :fetchedUsers="fetchedUsers"
            :sendWSMessage="sendWSMessage"
            :currentViewedRoom="currentViewedRoom"
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

import { WSMessageType, MessageType } from "./Constants";

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

    userID: store.state.email,
    newRoomName: "",

    showAddRoomDialog: false,
    showNotificationDialog: false,
    showChatPage: false,

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

        console.log(this.joinedRooms.length, this.indexOfCurrentViewedRoom);
        this.joinedRooms[
          this.indexOfCurrentViewedRoom
        ].roomIcon = this.currentViewedRoom.roomIcon;

        this.$nextTick(() => this.scrollToBottomOfChatPage());
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
          message: message,
          index: 0,
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
        this.updateRoomContentPage();
        this.currentViewedRoom.messages.push(message);
        this.$nextTick(() => this.scrollToBottomOfChatPage());
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
          time: "",
          name: "",
          userID: "",
          roomID: "",
          index: 0,
        });

        console.log("sent notification");
      }

      this.updateRecentMessagePreview(sentRequest.roomID, message);
    },

    updateRecentMessagePreview: function (roomID: string, message: string) {
      this.recentChatPreview[roomID] = message;
    },

    updateRoomContentPage: function () {
      this.$children[1].$forceUpdate();
    },

    updateChatRoomPage: function () {
      this.$children[2].$forceUpdate();
    },

    scrollToBottomOfChatPage: function () {
      const scrollHeight = this.$children[2].$el.querySelector("#messages");
      if (scrollHeight) scrollHeight.scrollTop = scrollHeight.scrollHeight;
    },

    changeNumberOfUnreadNotification: async function (
      roomID: string,
      addToUnreadNotifs: boolean
    ) {
      this.unreadRoomMessages[roomID] = addToUnreadNotifs;
      this.updateRoomContentPage();

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
      console.log("Received");

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
          console.log(jsonContent.userID, this.usersOnline);
          if (this.usersOnline[jsonContent.userID]) {
            this.usersOnline[jsonContent.userID].isOnline = jsonContent.status;
            console.log("updated");
            if (this.showChatPage) {
              console.log("final update");
              this.$nextTick(this.updateChatRoomPage);
            }
          }
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