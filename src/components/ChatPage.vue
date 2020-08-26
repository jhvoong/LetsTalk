<template>
  <v-row no-gutters style="height: 100vh;">
    <v-col cols="12" style="height: 10vh;">
      <v-app-bar tile flat height="100%" width="100%">
        <v-app-bar-nav-icon class="mx-5">
          <v-avatar>
            <v-img src="../assets/unilag.svg"></v-img>
          </v-avatar>
        </v-app-bar-nav-icon>

        <v-toolbar-title>
          <b>Room SEES2020</b>
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

    <v-col cols="12" style="height: 74vh;">
      <v-container class="overflow-y-auto scroll-behavior-smooth" style="height: 78vh;" fluid>
        <!-- <v-list width="70%">
          <v-list-item>
            <v-list-item-avatar height="min-content" width="min-content">
              <v-avatar rounded>
                <v-img height="100px" width="100px" contain :src="require('../assets/unilag.svg')"></v-img>
              </v-avatar>
            </v-list-item-avatar>

            <v-list-item-content>
              <v-card>

              </v-card>
            </v-list-item-content>
          </v-list-item>
        </v-list>-->

        <v-row>
          <v-col cols="12">
            <v-row style="width: 70%">
              <v-col cols="12">
                <v-row>
                  <v-avatar class="mx-2" style="align-self: flex-end;" rounded>
                    <v-img
                      height="50px"
                      width="50px"
                      contain
                      :src="require('../assets/unilag.svg')"
                    ></v-img>
                  </v-avatar>
                  <v-card flat :color="vuetify.framework.theme.dark?'':'#DCDCDC'" width="80%">
                    <v-card-text>Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.</v-card-text>
                  </v-card>
                </v-row>
              </v-col>
            </v-row>
          </v-col>

          <v-col cols="12">
            <div align="right">
              <v-card flat :color="vuetify.framework.theme.dark?'#778899':'#bbeaff'" width="60%">
                <v-card-text
                  align="left"
                >Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.</v-card-text>
              </v-card>
            </div>
          </v-col>

          <v-col cols="12">
            <div align="right">
              <v-card flat :color="vuetify.framework.theme.dark?'#778899':'#bbeaff'" width="60%">
                <v-card-text
                  align="left"
                >Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.</v-card-text>
              </v-card>
            </div>
          </v-col>

          <v-col cols="12">
            <v-row style="width: 70%">
              <v-col cols="12">
                <v-row>
                  <v-avatar class="mx-2" style="align-self: flex-end;" rounded>
                    <v-img
                      height="50px"
                      width="50px"
                      contain
                      :src="require('../assets/unilag.svg')"
                    ></v-img>
                  </v-avatar>
                  <v-card flat :color="vuetify.framework.theme.dark?'':'#DCDCDC'" width="80%">
                    <v-card-text>Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.</v-card-text>
                  </v-card>
                </v-row>
              </v-col>
            </v-row>
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
        dense
        no-resize
        rows="3"
        rounded
        clearable
        :readonly="disableTextField"
      ></v-textarea>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import Vue from "vue";
import vuetify from "@/plugins/vuetify";

export default Vue.extend({
  name: "ChatPage",
  data: () => ({
    vuetify: vuetify,
  }),
  methods: {},
});
</script>

<style>
</style>