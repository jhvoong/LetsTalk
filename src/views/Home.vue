<template>
  <v-row no-gutters class="fill-height">
    <v-col cols="1">
      <InnerSidebar />
    </v-col>
    <v-col cols="3">
      <OuterSidebar :joinedRooms="joinedRooms" />
    </v-col>
    <v-col cols="8">
      <ChatPage :send="sendMessage" />
    </v-col>
  </v-row>
</template>

<script lang="ts">
let socket: WebSocket;
// @ is an alias to /src
import Vue from "vue";
import store from "@/store";
import router from "@/router";

import { MessageType } from "./Constants";
// import { JoinedRoom } from "./Types";

import InnerSidebar from "../components/InnerSidebar.vue";
import OuterSidebar from "../components/OuterSidebar.vue";
import ChatPage from "../components/ChatPage.vue";
import { JoinedRoom } from "./Types";

export default Vue.extend({
  name: "Home",
  components: {
    InnerSidebar,
    OuterSidebar,
    ChatPage,
  },

  data: () => ({
    joinedRooms: new Array<JoinedRoom>(),
    roomJoinRequest: [],
  }),

  methods: {
    onWSMessage: function (event: MessageEvent) {
      const jsonContent = JSON.parse(event.data);

      switch (jsonContent.msgType) {
        case MessageType.UnauthorizedAccess:
          this.onUnAuthorizedAccess();
          break;

        case MessageType.WebsocketOpen:
          this.joinedRooms = jsonContent.joinedRooms;
          this.roomJoinRequest = jsonContent.roomJoinRequest;
          console.log(this.joinedRooms[0].roomID);
          break;
      }
    },

    onWSOpen: function () {
      console.log("Websocket open.");

      // Fetch users content from API.
      // Contents that are to be fetched from API are, Registered rooms and room request.
      const message = {
        msgType: MessageType.WebsocketOpen,
      };

      socket.send(JSON.stringify(message));
    },

    onWSError: function (event: Event) {
      // Reconnect Websocket if not UnAuthorized.
      console.log("Websocket errored.", event);
    },

    onWSClose: function (closeEvent: CloseEvent) {
      // We can create a dialog and ask users if they want to reconnect.
      console.log("Websocket closed.", closeEvent);
    },

    onUnAuthorizedAccess: function () {
      router.push("/login");
      socket.close();
    },

    sendMessage: function (message: string) {
      if (typeof message != "string") {
        return;
      }

      socket.send(message);
    },
  },

  created() {
    // Verify if token is specified.
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

    socket.onerror = this.onWSError;

    socket.onopen = this.onWSOpen;

    socket.onmessage = this.onWSMessage;

    socket.onclose = this.onWSClose;
  },
});
</script>

