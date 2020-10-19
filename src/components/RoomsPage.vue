<template>
  <v-row no-gutters style="justify-content: center">
    <v-col class="my-3 mx-5" cols="12" sm="10">
      <v-text-field
        rounded
        filled
        prepend-inner-icon="mdi-magnify"
        flat
        placeholder="Search for conversations"
      ></v-text-field>
    </v-col>

    <v-col cols="12" class="text-center">
      <v-btn
        :outlined="currentRoomFilter == joinedRoomFilter.All"
        @click="currentRoomFilter = joinedRoomFilter.All"
        text
        rounded
        >{{ joinedRoomFilter.All }}</v-btn
      >

      <v-btn
        :outlined="currentRoomFilter == joinedRoomFilter.Read"
        @click="currentRoomFilter = joinedRoomFilter.Read"
        text
        rounded
        >{{ joinedRoomFilter.Read }}</v-btn
      >

      <v-btn
        :outlined="currentRoomFilter == joinedRoomFilter.Unread"
        @click="currentRoomFilter = joinedRoomFilter.Unread"
        text
        rounded
        >{{ joinedRoomFilter.Unread }}</v-btn
      >
    </v-col>

    <v-col cols="12">
      <v-subheader>
        <b>Discussions</b>
      </v-subheader>
    </v-col>

    <v-container style="height: 70vh" class="overflow-y-auto" cols="12">
      <v-list nav tile dense three-line>
        <v-list-item-group v-model="roomIndex" mandatory>
          <v-list-item
            v-for="(joinedRoom, index) in filteredRooms"
            :key="index"
            @click="loadChatContent(joinedRoom.roomID)"
          >
            <v-list-item-avatar>
              <v-badge
                bordered
                overlap
                bottom
                dot
                :value="
                  unreadRoomMessages[joinedRoom.roomID]
                    ? unreadRoomMessages[joinedRoom.roomID]
                    : false
                "
                offset-x="10"
                offset-y="10"
              >
                <v-avatar size="30">
                  <v-img
                    :src="
                      joinedRoom.roomIcon
                        ? joinedRoom.icon
                        : require('../assets/unilag.svg')
                    "
                  ></v-img>
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
                    >{{ joinedRoom.roomName }}</v-col
                  >
                </v-row>
              </v-list-item-title>

              <v-list-item-subtitle class="d-inline-block text-truncate">{{
                recentChatPreview[joinedRoom.roomID]
                  ? recentChatPreview[joinedRoom.roomID]
                  : ""
              }}</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </v-container>
  </v-row>
</template>

<script lang="ts">
import Vue from "vue";
import { Prop } from "vue/types/options";
import store from "@/store";

import { JoinedRoom, RecentChatPreview, UnreadRooms } from "../views/Types";
import { WSMessageType, JoinedRoomsFilter } from "../views/Constants";

export default Vue.extend({
  name: "RoomsPage",

  props: {
    joinedRooms: Array as Prop<JoinedRoom[]>,
    recentChatPreview: {} as Prop<RecentChatPreview>,
    unreadRoomMessages: {} as Prop<UnreadRooms>,

    sendWSMessage: Function,
    changeViewedRoomIndex: Function,

    indexOfCurrentViewedRoom: Number,
  },

  data: () => ({
    userID: store.state.email,
    joinedRoomFilter: JoinedRoomsFilter,
    currentRoomFilter: JoinedRoomsFilter.All,

    roomIndex: 0,
    filteredRooms: [] as JoinedRoom[],
  }),

  methods: {
    loadChatContent: function (roomID: string) {
      const message = {
        msgType: WSMessageType.RequestMessages,
        userID: this.userID,
        roomID: roomID,
        firstLoad: true,
      };

      this.sendWSMessage(JSON.stringify(message));
    },

    filterRooms: function (retrieveOnlyUnread: boolean): JoinedRoom[] {
      const joinedRooms: JoinedRoom[] = [];

      this.joinedRooms.forEach((joinedRoom: JoinedRoom) => {
        let isUnread = this.unreadRoomMessages[joinedRoom.roomID];
        if (!isUnread) {
          isUnread = false;
        }
        if (isUnread == retrieveOnlyUnread) {
          joinedRooms.push(joinedRoom);
        }
      });

      return joinedRooms;
    },

    roomFilterWatcher: function () {
      switch (this.currentRoomFilter) {
        case this.joinedRoomFilter.Unread:
          this.filteredRooms = this.filterRooms(true);
          return;
        case this.joinedRoomFilter.Read:
          this.filteredRooms = this.filterRooms(false);
          return;
      }

      console.log("all was selected");
      this.filteredRooms = this.joinedRooms;
    },

    getRoomIndex: function (
      roomID: string,
      joinedRooms: JoinedRoom[]
    ): Promise<number> {
      return new Promise((): number => {
        for (let i = 0; i < joinedRooms.length; i++) {
          if (joinedRooms[0].roomID === roomID) {
            return i;
          }
        }

        throw new Error("error, roomID not found");
      });
    },
  },

  watch: {
    indexOfCurrentViewedRoom: function () {
      // Change room index of filtered room if indexOfCurrentViewedRoom changes.
      // indexOfCurrentViewedRoom tracks index of unfiltered room.
      this.getRoomIndex(
        this.joinedRooms[this.indexOfCurrentViewedRoom].roomID,
        this.filteredRooms
      )
        .then((index: number) => {
          this.roomIndex = index;
        })
        .catch(() => {
          // This is assumed there is no rooms.
          this.roomIndex = 0;
        });
    },

    roomIndex: function () {
      console.log("roomIndex");
      this.getRoomIndex(
        this.joinedRooms[this.indexOfCurrentViewedRoom].roomID,
        this.joinedRooms
      )
        .then((index: number) => {
          this.changeViewedRoomIndex(index);
        })
        .catch(() => {
          // If room index is not allocated, we change to first index.
          this.changeViewedRoomIndex(0);
        });
    },

    joinedRooms: function () {
      console.log("joinedRooms");
      this.roomFilterWatcher();
    },

    currentRoomFilter: function () {
      console.log("currentRoomFilter");
      this.roomFilterWatcher();
    },
  },

  created() {
    this.filteredRooms = this.joinedRooms;
  },
});
</script>

<style>
.button_default {
  font-family: sans-serif;
  font-weight: 600;
  font-size: small;
}
</style>