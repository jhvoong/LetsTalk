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

        <v-menu
          v-model="showAddUsersDialog"
          left
          rounded="r-xl"
          open-on-hover
          :close-on-content-click="false"
          append-icon="mdi-magnify"
          @keyup.enter.exact="searchUsers"
          @click:append="searchUsers"
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
        <v-row dense>
          <v-col cols="12" v-for="(message,index) in currentViewedRoom.messages" :key="index">
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
                <v-avatar style="align-self: flex-end;" rounded>
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
"use strict";

import Vue from "vue";
import vuetify from "@/plugins/vuetify";
import { Prop } from "vue/types/options";
import store from "@/store";

import { RoomPageDetails, FetchedUsers } from "../views/Types";
import { WSMessageType, MessageType } from "../views/Constants";

export default Vue.extend({
  name: "ChatPage",

  props: {
    currentViewedRoom: {} as Prop<RoomPageDetails>,
    fetchedUsers: Array as Prop<FetchedUsers[]>,

    sendWSMessage: Function,
    clearFetchedUsers: Function,
  },

  data: () => ({
    vuetify: vuetify,
    messageContent: "",
    newRoomName: "",
    searchUserName: "",

    showAddUsersDialog: false,

    selectedUsersToAddToRoom: [] as string[],

    messageType: MessageType,
    userID: store.state.email,
    userName: store.state.name,
  }),

  methods: {
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

    loadMoreMessages: function () {
      const message = {
        msgType: WSMessageType.RequestMessages,
        userID: this.userID,
        roomID: this.currentViewedRoom.roomID,

        messageCount: this.currentViewedRoom.messages[
          this.currentViewedRoom.messages.length - 1
        ].index,
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
  },
});
</script>

<style>
</style>