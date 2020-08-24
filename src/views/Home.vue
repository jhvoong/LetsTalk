<template>
  <div class="home" align="center">
    <img alt="Unilag logo" src="../assets/unilag.svg" />
    <p class="text-center">Welcome to the university of lagos chatting platform</p>
  </div>
</template>

<script lang="ts">
// @ is an alias to /src
import Vue from "vue";
import store from "@/store";
import router from "@/router";
import * as constants from "./Constants";

let socket: WebSocket;

export default Vue.extend({
  name: "Home",
  components: {},

  data: () => ({}),

  methods: {
    onWSMessage: function (event: MessageEvent) {
      const jsonContent: any = JSON.parse(event.data);

      switch (jsonContent.msgType) {
        case constants.MessageType.UnauthorizedAccess:
          this.onUnAuthorizedAccess();
          break;
      }
    },

    onWSOpen: function () {
      console.log("Websocket open.");
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
      console.log("Url not defined");
      router.push("/login");
      return;
    }

    console.log(URLs);
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

