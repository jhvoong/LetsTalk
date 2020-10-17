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
        z-index: 0;
        object-fit: cover;
      "
      id="videoID"
      autoplay
    ></video>

    <v-fab-transition>
      <v-row
        style="
          width: 100vw;
          height: 100vh;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          position: relative;
          z-index: 1;
        "
        justify="center"
        v-if="showVideoOptions"
      >
        <v-col cols="12" class="mx-5 d-flex justify-end">
          <v-btn fab dark small bottom left v-if="showVideoOptions">
            <v-icon>mdi-account-circle</v-icon>
          </v-btn>
        </v-col>

        <v-col cols="12" class="my-10 d-flex align-end">
          <v-row justify="center" v-if="showVideoOptions">
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

            <v-btn
              class="ml-2"
              color="red"
              small
              fab
              dark
              @click="endCallSession()"
            >
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
              @click="changeFullScreenStatus()"
            >
              <v-icon>mdi-fullscreen</v-icon>
            </v-btn>
          </v-row>
        </v-col>
      </v-row>
    </v-fab-transition>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
export default Vue.extend({
  name: "CallUI",

  props: {
    isCallPublisher: Boolean,

    changeVideo: Function,
    startDesktopSharing: Function,
    stopDesktopSharing: Function,
    endCallSession: Function,
    changeAudio: Function,
  },

  data: () => ({
    showVideoOptions: false,
    isDesktopShared: false,
    videoOn: false,
    audioOn: false,
    fullScreen: false,

    videoIcon: "mdi-video",
    audioIcon: "mdi-microphone",
    screenMode: "mdi-fullscreen",
  }),

  methods: {
    desktopShare: function () {
      if (!this.isCallPublisher) {
        return;
      }

      console.log("Starting desktop share");
      if (!this.isDesktopShared) {
        this.startDesktopSharing();
      } else {
        this.stopDesktopSharing();
      }

      this.isDesktopShared = !this.isDesktopShared;
    },

    changeVideoStatus: function () {
      if (!this.isCallPublisher || this.isDesktopShared) {
        return;
      }

      this.videoOn = !this.videoOn;

      if (this.videoOn) {
        this.videoIcon = "mdi-video";
      } else {
        this.videoIcon = "mdi-video-off";
      }

      this.changeVideo();
    },

    changeAudioStatus: function () {
      if (this.audioOn) {
        this.audioIcon = "mdi-microphone-off";
      } else {
        this.audioIcon = "mdi-microphone";
      }

      this.audioOn = !this.audioOn;
      this.changeAudio(this.audioOn);
    },

    changeFullScreenStatus: function () {
      if (!document.fullscreenElement) {
        document.documentElement.requestFullscreen();
      } else if (document.exitFullscreen) {
        document.exitFullscreen();
      }

      this.fullScreen = !this.fullScreen;
    },
  },
});
</script>