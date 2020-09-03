<template>
  <v-row no-gutters style="height: 100vh;">
    <v-col cols="12" style="height: 10vh;">
      <v-app-bar tile flat height="100%" width="100%">
        <v-app-bar-nav-icon class="mx-5">
          <v-avatar>
            <v-img :src="currentViewedRoom.roomIcon"></v-img>
          </v-avatar>
        </v-app-bar-nav-icon>

        <v-toolbar-title>
          <b>{{currentViewedRoom.roomName}}</b>
        </v-toolbar-title>

        <v-spacer></v-spacer>

        <v-btn class="mx-2" @click="startSession()" icon x-large>
          <v-icon>mdi-account-multiple-plus-outline</v-icon>
        </v-btn>
        <v-btn class="mx-2" @click="startSession()" icon x-large>
          <v-icon>mdi-video</v-icon>
        </v-btn>
        <v-btn class="mx-2" @click="startSession()" icon x-large>
          <v-icon>mdi-information</v-icon>
        </v-btn>
      </v-app-bar>
    </v-col>

    <v-col cols="12" style="height: 75vh;">
      <v-container class="overflow-y-auto scroll-behavior-smooth" style="height: 78vh;" fluid>
        <v-row>
          <v-col cols="12" v-for="(message,index) in currentViewedRoom.messages" :key="index">
            <template v-if="message.type===messageType.Message">
              <v-row v-if="message.userID===userID">
                <div align="right">
                  <v-card
                    flat
                    :color="vuetify.framework.theme.dark?'#778899':'#bbeaff'"
                    width="60%"
                  >
                    <v-card-text align="left">{{message.message}}</v-card-text>

                    <v-card-subtitle align="right" class="ml-auto">
                      <b>{{message.userID}} {{message.time}}</b>
                    </v-card-subtitle>
                  </v-card>
                </div>
              </v-row>

              <v-row v-else>
                <v-avatar style="align-self: flex-end;" rounded>
                  <v-img height="50px" width="50px" contain :src="require('../assets/unilag.svg')"></v-img>
                </v-avatar>
                <v-card flat :color="vuetify.framework.theme.dark?'':'#DCDCDC'" width="70%">
                  <v-card-text>{{message.message}}</v-card-text>
                  <v-card-subtitle align="right" class="ml-auto">
                    <b>{{message.userID}} {{message.time}}</b>
                  </v-card-subtitle>
                </v-card>
              </v-row>
            </template>

            <template v-else>
              <v-col v-if="message.type===messageType.File" cols="12">
                <div align="center">
                  <v-chip href="https://github.com/metaclips">
                    <b>{{message.message}} sent by {{message.name}}. Click to download.</b>
                  </v-chip>
                </div>
              </v-col>

              <v-col v-if="message.type===messageType.Info" cols="12">
                <div align="center">
                  <v-chip href="https://github.com/metaclips">
                    <b>{{message.message}}.</b>
                  </v-chip>
                </div>
              </v-col>

              <v-col v-if="message.type===messageType.ClassSessionLink" cols="12">
                <div align="center">
                  <v-chip href="https://github.com/metaclips">
                    <b>Class session started by {{message.name}}. Click to join.</b>
                  </v-chip>
                </div>
              </v-col>

              <v-col v-if="message.type===messageType.ClassSession" cols="12">
                <div align="center">
                  <v-chip href="https://github.com/metaclips">
                    <b>Class session recording by {{message.name}}. Click to download.</b>
                  </v-chip>
                </div>
              </v-col>
            </template>
          </v-col>
        </v-row>
      </v-container>
    </v-col>

    <v-col cols="12" style="height: 10vh;">
      <v-textarea
        class="mx-5"
        v-model="messageContent"
        prepend-inner-icon="mdi-emoticon"
        prepend-icon="mdi-paperclip"
        append-outer-icon="mdi-send"
        solo
        hide-details="auto"
        no-resize
        rows="3"
        rounded
        clearable
        @click:append-outer="sendMessage"
        @keyup.enter.exact="sendMessage"
      ></v-textarea>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import Vue from "vue";
import vuetify from "@/plugins/vuetify";
import { Prop } from "vue/types/options";
import store from "@/store";

import { RoomPageDetails } from "../views/Types";
import { WSMessageType, MessageType } from "../views/Constants";

export default Vue.extend({
  name: "ChatPage",

  props: {
    currentViewedRoom: {} as Prop<RoomPageDetails>,

    sendWSMessage: Function,
  },

  data: () => ({
    vuetify: vuetify,
    messageContent: "",
    newRoomName: "",

    messageType: MessageType,

    userID: store.state.email,
  }),

  methods: {
    sendMessage: function () {
      if (!this.messageContent.match(/\S/)) return;

      const message = {
        msgType: WSMessageType.NewMessage,
        message: this.messageContent,
        userID: this.userID,
        type: this.messageType.Message,
      };

      this.sendWSMessage(JSON.stringify(message));
      this.messageContent = "";
    },

    loadMoreMessages: function () {
      const message = {
        msgType: WSMessageType.RequestMessages,
        roomID: this.currentViewedRoom.roomID,

        messageCount: this.currentViewedRoom.messages[
          this.currentViewedRoom.messages.length - 1
        ].index,
      };

      this.sendWSMessage(JSON.stringify(message));
    },
  },
});
</script>

<style>
</style>