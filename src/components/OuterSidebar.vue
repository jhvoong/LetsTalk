<template>
  <v-container fluid>
    <v-row style="height: 25vh; justify-content: center;">
      <v-col cols="12" sm="10">
        <v-text-field
          rounded
          filled
          prepend-inner-icon="mdi-magnify"
          flat
          placeholder="Search for conversations"
        ></v-text-field>
      </v-col>

      <v-col cols="12" class="text-center">
        <v-btn outlined text rounded>All</v-btn>
        <v-btn text rounded>Read</v-btn>
        <v-btn text rounded>Unread</v-btn>
      </v-col>

      <v-col cols="12">
        <h2>Discussions</h2>
      </v-col>
    </v-row>

    <v-container style=" height: 72vh;" class="overflow-y-auto" cols="12">
      <v-list nav tile dense three-line>
        <v-list-item-group mandatory>
          <v-list-item
            v-for="(joinedRoom,index) in joinedRooms"
            :key="index"
            @click="loadChatContent(joinedRoom.roomID)"
          >
            <v-list-item-avatar>
              <v-badge bordered bottom color="deep-purple accent-4" dot offset-x="10" offset-y="10">
                <v-avatar size="30">
                  <v-img :src="joinedRoom.icon?joinedRoom.icon:require('../assets/unilag.svg')"></v-img>
                </v-avatar>
              </v-badge>
            </v-list-item-avatar>

            <v-list-item-content>
              <v-list-item-title>
                <v-row no-gutters>
                  <v-col
                    cols="auto"
                    class="d-inline-block text-truncate"
                    style="width: 70%"
                    justify="start"
                  >{{joinedRoom.roomName}}</v-col>
                </v-row>
              </v-list-item-title>

              <v-list-item-subtitle
                class="d-inline-block text-truncate"
              >This is text with nonesensejsjsjjsjsjsjsjsjsjkskskkskskksksk heads message</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </v-container>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import { Prop } from "vue/types/options";

import { JoinedRoom, RecentChatPreview } from "../views/Types";

export default Vue.extend({
  name: "OuterSidebar",

  props: {
    joinedRooms: Array as Prop<JoinedRoom[]>,
    recentChatPreview: {} as Prop<RecentChatPreview>,

    loadChatContent: Function,
  },

  data: () => ({}),

  methods: {},
});
</script>

<style>
.button_default {
  font-family: sans-serif;
  font-weight: 600;
  font-size: small;
}
</style>