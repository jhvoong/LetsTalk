(%define "callUI"%)
<template>
  <v-container
    fill-height
    align="center"
    justify="center"
    style="min-width: 100%; min-height: 100%; background-color: black"
    @click="showVideoOptions = !showVideoOptions"
  >
    <video
      style="
        width: 100%;
        height: 100%;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        position: absolute;
        object-fit: cover;
      "
      id="videoID"
      autoplay
    ></video>

    <v-expand-transition>
      <v-row
        justify="center"
        v-if="showVideoOptions"
        style="
          bottom: 1rem;
          padding: 20px;
          left: 50%;
          margin-left: -150px;
          position: fixed;
          z-index: 2;
        "
      >
        <v-btn
          class="ml-2"
          small
          fab
          dark
          @click="desktopShare()"
          :color="isDesktopShared ? 'blue' : ''"
        >
          <v-icon>mdi-desktop-mac</v-icon>
        </v-btn>

        <v-btn class="ml-2" fab small dark @click="changeVideoStatus()">
          <v-icon>{{ videoIcon }}</v-icon>
        </v-btn>

        <v-btn class="ml-2" color="red" small fab dark @click="endCall()">
          <v-icon>mdi-phone-hangup</v-icon>
        </v-btn>

        <v-btn class="ml-2" small fab dark @click="changeAudioStatus()">
          <v-icon>{{ audioIcon }}</v-icon>
        </v-btn>

        <v-btn
          :color="fullScreen ? 'blue' : ''"
          class="ml-2"
          small
          fab
          dark
          @click="requestFullScreen()"
        >
          <v-icon>mdi-fullscreen</v-icon>
        </v-btn>
      </v-row>
    </v-expand-transition>

    <v-dialog persistent v-model="onClassSessionError" width="300px">
      <v-card class="d-inline-block mx-auto">
        <v-card-text class="text--secondary" style="height: 250px">
          <v-row>
            <v-col cols="12"></v-col>
            <v-col cols="12">{{ class_session_error }}</v-col>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="green darken-1" text @click="closeDialog"
            >Close Dialog</v-btn
          >
          <v-btn color="red darken-1" text @click="endCall">End Call</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>
(%end%)