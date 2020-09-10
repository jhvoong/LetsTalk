<template>
  <v-row no-gutters style="height: 100vh;">
    <v-col cols="12" style="height: 10vh;">
      <v-app-bar tile flat height="100%" width="100%">
        <v-app-bar-nav-icon class="mx-5">
          <v-avatar>
            <v-img
              :src="currentViewedRoom.roomIcon?currentViewedRoom.roomIcon:require('../assets/unilag.svg')"
            ></v-img>
          </v-avatar>
        </v-app-bar-nav-icon>

        <v-toolbar-title>
          <b>{{currentViewedRoom.roomName}}</b>
        </v-toolbar-title>

        <v-spacer></v-spacer>

        <v-menu
          v-model="showAddUsersDialog"
          left
          rounded="r-xl"
          :close-on-content-click="false"
          nudge-width="400"
          offset-y
        >
          <template v-slot:activator="{ on, attrs }">
            <v-btn class="mx-2" v-bind="attrs" v-on="on" icon x-large>
              <v-icon>mdi-account-multiple-plus-outline</v-icon>
            </v-btn>
          </template>

          <v-card align="center">
            <v-row style="max-width: 40vh;" align="center" justify="center">
              <v-col cols="12"></v-col>
              <v-col cols="12" sm="8">
                <v-text-field
                  v-model="searchUserName"
                  append-icon="mdi-magnify"
                  @keyup.enter.exact="searchUsers"
                  @click:append="searchUsers"
                  rounded
                  filled
                  placeholder="Search for contact email"
                />
              </v-col>

              <v-col cols="12">
                <v-container
                  fluid
                  class="ml-auto overflow-y-auto scroll-behavior-smooth"
                  style="max-height: 40vh;"
                >
                  <v-checkbox
                    class="text--truncate"
                    dense
                    multiple
                    v-for="(user, i) in fetchedUsers"
                    :key="i"
                    :label="user.name+' ['+user.userID+']'"
                    :value="user.userID"
                    v-model="selectedUsersToAddToRoom"
                  ></v-checkbox>
                </v-container>
              </v-col>
            </v-row>

            <v-card-actions>
              <v-btn @click="closeAddUserDialog" color="red" text>exit</v-btn>
              <v-btn
                @click="requestUserToJoinRoom"
                :disabled="selectedUsersToAddToRoom.length===0"
                color="green"
                text
              >add users</v-btn>
            </v-card-actions>
          </v-card>
        </v-menu>

        <v-btn class="mx-2" icon x-large>
          <v-icon>mdi-video</v-icon>
        </v-btn>

        <v-menu
          v-model="showOnlineUsersDialog"
          left
          :close-on-content-click="false"
          rounded="r-xl"
          nudge-width="400"
          offset-y
        >
          <template v-slot:activator="{ on, attrs }">
            <v-btn class="mx-2" v-bind="attrs" v-on="on" icon x-large>
              <v-icon>mdi-information</v-icon>
            </v-btn>
          </template>

          <v-card>
            <v-row class="mx-10">
              <v-col cols="12"></v-col>
              <v-col
                cols="12"
                v-for="(onlineUser,index) in currentViewedRoom.registeredUsers"
                :key="index"
              >
                <template v-if="onlineUser==userID">
                  <v-badge inline dot color="green"></v-badge>
                  <span class="mx-4">
                    <b>You</b>
                  </span>
                </template>

                <template v-else>
                  <v-badge
                    inline
                    dot
                    :color="associates[onlineUser]&&associates[onlineUser].isOnline===true ? 'green' : 'red'"
                  ></v-badge>
                  <span class="mx-4">
                    <b>{{associates[onlineUser]?associates[onlineUser].name:""}} [{{onlineUser}}]</b>
                  </span>
                </template>
              </v-col>
            </v-row>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn @click="exitRoom" color="red" text>exit room</v-btn>
              <v-spacer></v-spacer>
            </v-card-actions>
          </v-card>
        </v-menu>
      </v-app-bar>
    </v-col>

    <v-col cols="12" style="height: 75vh;">
      <v-container
        id="messages"
        class="overflow-y-auto scroll-behavior-smooth"
        style="height: 75vh;"
        fluid
      >
        <v-row v-scroll:#messages="onScroll" dense>
          <v-col
            class="my-2"
            cols="12"
            v-for="(message,index) in currentViewedRoom.messages"
            :key="index"
          >
            <template v-if="message.type===messageType.Message">
              <v-row dense align="end" justify="end" v-if="message.userID===userID">
                <v-card
                  flat
                  :color="vuetify.framework.theme.dark?'#778899':'#bbeaff'"
                  max-width="65%"
                >
                  <v-card-text class="text--wrap">
                    <b>{{message.message}}</b>
                  </v-card-text>

                  <v-card-subtitle align="right">{{message.userID}} {{message.time}}</v-card-subtitle>
                </v-card>
              </v-row>

              <v-row dense v-else>
                <v-avatar class="mx-2" style="align-self: flex-end;" rounded>
                  <v-img height="50px" width="50px" contain :src="require('../assets/unilag.svg')"></v-img>
                </v-avatar>
                <v-card flat :color="vuetify.framework.theme.dark?'':'#DCDCDC'" max-width="65%">
                  <v-card-text>
                    <b>{{message.message}}</b>
                  </v-card-text>
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
                  <v-chip>
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
      <v-expand-transition>
        <v-form v-model="fileInputValid" v-show="showFileInput">
          <v-file-input
            :rules="rules"
            v-show="showFileInput"
            class="mx-5"
            v-model="file"
            counter
            label="File input"
            placeholder="Select your file"
            outlined
            append-outer-icon="mdi-send"
            :show-size="1000"
            @click:append-outer="sendFile"
          ></v-file-input>
        </v-form>
      </v-expand-transition>

      <v-expand-transition>
        <v-textarea
          class="mx-5"
          v-show="showTextField"
          v-model="messageContent"
          prepend-inner-icon="mdi-emoticon"
          prepend-icon="mdi-paperclip"
          append-outer-icon="mdi-send"
          solo
          hide-details="auto"
          clearable
          no-resize
          rows="3"
          rounded
          @click:prepend="hideTextField"
          @click:append-outer="sendMessage"
          @keyup.enter.exact="sendMessage"
        ></v-textarea>
      </v-expand-transition>
    </v-col>
  </v-row>
</template>

<script lang="ts">
"use strict";

import Vue from "vue";
import vuetify from "@/plugins/vuetify";
import { Prop } from "vue/types/options";
import store from "@/store";

import { RoomPageDetails, FetchedUsers, UsersOnline } from "../views/Types";
import { WSMessageType, MessageType } from "../views/Constants";

export default Vue.extend({
  name: "ChatPage",

  props: {
    currentViewedRoom: {} as Prop<RoomPageDetails>,
    fetchedUsers: Array as Prop<FetchedUsers[]>,
    associates: {} as Prop<UsersOnline>,

    sendWSMessage: Function,
    clearFetchedUsers: Function,
  },

  data: () => ({
    vuetify: vuetify,
    messageContent: "",
    newRoomName: "",
    searchUserName: "",

    showAddUsersDialog: false,
    showOnlineUsersDialog: false,
    showTextField: true,
    showFileInput: false,
    fileUploadOrDownload: false,
    fileInputValid: false,

    selectedUsersToAddToRoom: [] as string[],
    file: {} as File,

    messageType: MessageType,
    userID: store.state.email,
    userName: store.state.name,

    fileProgress: 0,

    rules: [
      (value: File) => !!value || "File cannot be empty",
      (value: File) =>
        !value || value.size < 256 * 1024 || "File size limit is 256MB!",
    ],
  }),

  methods: {
    onScroll: function (e: { target: Element }) {
      const messages = this.currentViewedRoom.messages;
      if (
        e.target &&
        e.target.scrollTop < 100 &&
        messages.length > 0 &&
        messages[0].index > 1
      ) {
        this.loadMoreMessages();
      }
    },

    sendMessage: function () {
      if (!this.messageContent.match(/\S/)) return;

      const message = {
        msgType: WSMessageType.NewMessage,
        userID: this.userID,
        roomID: this.currentViewedRoom.roomID,
        message: this.messageContent,
        type: this.messageType.Message,
      };

      this.sendWSMessage(JSON.stringify(message));
      this.messageContent = "";
    },

    sendFile: function () {
      console.log(this.fileInputValid);
      if (this.fileInputValid == false) return;

      const reader = new FileReader();
      reader.onloadend = function (event: ProgressEvent<FileReader>) {
        if (event.target) console.log(event.target.result);
      };

      reader.readAsDataURL(this.file);
    },

    loadMoreMessages: function () {
      const message = {
        msgType: WSMessageType.RequestMessages,
        userID: this.userID,
        roomID: this.currentViewedRoom.roomID,

        messageCount: this.currentViewedRoom.messages[0].index - 1,
      };

      this.sendWSMessage(JSON.stringify(message));
    },

    searchUsers: function () {
      const message = {
        msgType: WSMessageType.SearchUser,
        userID: this.userID,
        searchText: this.searchUserName,
      };

      this.sendWSMessage(JSON.stringify(message));
    },

    requestUserToJoinRoom: function () {
      const message = {
        msgType: WSMessageType.RequestUsersToJoinRoom,
        userID: this.userID,
        roomID: this.currentViewedRoom.roomID,
        roomName: this.currentViewedRoom.roomName,
        requestingUserName: this.userName,
        requestingUserID: this.userID,
        users: this.selectedUsersToAddToRoom,
      };

      this.sendWSMessage(JSON.stringify(message));
      this.closeAddUserDialog();
    },

    closeAddUserDialog: function () {
      this.selectedUsersToAddToRoom = [];
      this.showAddUsersDialog = false;
      this.clearFetchedUsers();
    },

    exitRoom: function () {
      const message = {
        msgType: WSMessageType.ExitRoom,
        roomID: this.currentViewedRoom.roomID,
        userID: this.userID,
      };

      this.sendWSMessage(JSON.stringify(message));
      this.showOnlineUsersDialog = false;
    },

    hideTextField: function () {
      this.showFileInput = true;
      this.showTextField = false;
    },

    revealTextField: function () {
      this.showFileInput = false;
      this.showTextField = true;
    },
  },
  watch: {
    file: function () {
      if (!this.file) {
        this.revealTextField();
      }
    },
  },
});
</script>
