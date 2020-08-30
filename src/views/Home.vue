<template>
  <v-row no-gutters class="fill-height">
    <v-col cols="1">
      <InnerSidebar />
    </v-col>
    <v-col cols="3">
      <OuterSidebar :joinedRooms="joinedRooms" :loadChatContent="loadChatContent" />
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
import { JoinedRoom, RoomPageDetails } from "./Types";

export default Vue.extend({
  name: "Home",
  components: {
    InnerSidebar,
    OuterSidebar,
    ChatPage,
  },

  data: () => ({
    joinedRooms: [] as JoinedRoom[],

    messages: {} as RoomPageDetails,
  }),

  methods: {
    onUnAuthorizedAccess: function () {
      router.push("/login");
      socket.close();
    },

    sendMessage: function (message: string) {
      socket.send(message);
    },

    loadChatContent: function (roomID: string) {
      const message = {
        msgType: MessageType.RequestMessages,
        roomID: roomID,
      };

      socket.send(JSON.stringify(message));
    },
  },

  created() {
    socket.onmessage = (event: MessageEvent) => {
      const jsonContent = JSON.parse(event.data);

      switch (jsonContent.msgType) {
        case MessageType.UnauthorizedAccess:
          this.onUnAuthorizedAccess();
          break;

        case MessageType.WebsocketOpen:
          this.joinedRooms = jsonContent.joinedRooms;
          break;

        case MessageType.RequestMessages:
          break;
      }
    };

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
        msgType: MessageType.WebsocketOpen,
      };

      socket.send(JSON.stringify(message));
    };

    socket.onclose = function (closeEvent: CloseEvent) {
      // We can create a dialog and ask users if they want to reconnect.
      console.log("Websocket closed.", closeEvent);
    };
  },
});
</script>

