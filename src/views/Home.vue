<template>
  <v-row no-gutters class="fill-height">
    <v-col cols="1">
      <InnerSidebar
        :activateAddRoomDialog="activateAddRoomDialog"
        :activateNotificationDialog="activateNotificationDialog"
        :deactivateAllDialogs="deactivateAllDialogs"
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
                    <v-card-text>{{joinRequest.requestingUserName}} [{{joinRequest.requestingUserID}}] wants you to join room {{joinRequest.RoomID}}</v-card-text>
                    <v-card-action>
                      <v-btn
                        @click="joinRoom(joinRequest.roomID, joinRequest.roomName, index)"
                        text
                        color="green"
                      >Join</v-btn>
                      <v-btn text color="red">Reject</v-btn>
                    </v-card-action>
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
        <v-col cols="4">
          <OuterSidebar :joinedRooms="joinedRooms" :sendWSMessage="sendWSMessage" />
        </v-col>

        <v-col cols="8">
          <ChatPage
            v-if="showChatPage"
            :sendWSMessage="sendWSMessage"
            :currentViewedRoom="currentViewedRoom"
          />
        </v-col>
      </v-row>
    </v-col>
  </v-row>
</template>

<script lang="ts">
let socket: WebSocket;
// @ is an alias to /src
import Vue from "vue";
import store from "@/store";
import router from "@/router";

import { WSMessageType } from "./Constants";

import InnerSidebar from "../components/InnerSidebar.vue";
import OuterSidebar from "../components/OuterSidebar.vue";
import ChatPage from "../components/ChatPage.vue";
import { JoinedRoom, RoomPageDetails, Messages, JoinRequest } from "./Types";

export default Vue.extend({
  name: "Home",
  components: {
    InnerSidebar,
    OuterSidebar,
    ChatPage,
  },

  data: () => ({
    joinedRooms: [] as JoinedRoom[],
    currentViewedRoom: {} as RoomPageDetails,
    joinRequests: [] as JoinRequest[],

    userID: store.state.email,
    newRoomName: "",

    showAddRoomDialog: false,
    showNotificationDialog: false,
    showChatPage: false,
  }),

  methods: {
    onUnAuthorizedAccess: function () {
      router.push("/login");
      socket.close();
    },

    onRequestMessages: function (roomDetails: RoomPageDetails) {
      if (roomDetails.firstLoad) {
        this.currentViewedRoom = roomDetails;
        if (this.currentViewedRoom.roomIcon == "") {
          this.currentViewedRoom.roomIcon = require("../assets/unilag.svg");
        }

        return;
      }

      roomDetails.messages.map((messages: Messages) => {
        this.currentViewedRoom.messages.unshift(messages);
      });
    },

    joinRoom: function (roomID: string, roomName: string, index: number) {
      const message = {
        joined: true,
        roomID: roomID,
        roomName: roomName,
        msgType: WSMessageType.JoinRoom,
        userID: this.userID,
      };

      socket.send(JSON.stringify(message));
      this.joinRequests.splice(index, 1);
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

  created() {
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

    socket.onmessage = (event: MessageEvent) => {
      const jsonContent = JSON.parse(event.data);

      switch (jsonContent.msgType) {
        case WSMessageType.UnauthorizedAccess:
          this.onUnAuthorizedAccess();
          break;

        case WSMessageType.WebsocketOpen:
          this.joinedRooms = jsonContent.joinedRooms;
          break;

        case WSMessageType.RequestMessages:
          this.showChatPage = true;
          this.onRequestMessages(jsonContent.roomPageDetails);
          break;

        case WSMessageType.JoinRoom:
          this.joinedRooms.unshift({
            roomID: jsonContent.roomID,
            roomIcon: "",
            roomName: jsonContent.roomName,
          });
          break;
      }
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