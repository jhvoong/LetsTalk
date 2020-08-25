<template>
  <v-card flat height="100%" align="center" outlined width="max-content">
    <v-row class="fill-height">
      <v-col class="my-2" cols="12">
        <v-avatar>
          <img src="../assets/unilag.svg" alt="Unilag Logo" />
        </v-avatar>
      </v-col>

      <v-col cols="12">
        <v-row>
          <v-col class="my-2" cols="12">
            <v-btn @click="onContact" x-large icon>
              <v-icon
                :color="highlightedSidebarOption=='Contact'?'blue':''"
                x-large
              >mdi-account-circle</v-icon>
            </v-btn>
          </v-col>

          <v-col class="my-2" cols="12">
            <v-btn @click="onMessage" x-large icon>
              <v-icon
                :color="highlightedSidebarOption=='Message'?'blue':''"
                x-large
              >mdi-message-outline</v-icon>
            </v-btn>
          </v-col>

          <v-col class="my-2" cols="12">
            <v-btn @click="onNotification" x-large icon>
              <v-icon :color="highlightedSidebarOption=='Notif'?'blue':''" x-large>mdi-bell-outline</v-icon>
            </v-btn>
          </v-col>
        </v-row>
      </v-col>

      <v-col class="mt-auto" cols="12">
        <v-row>
          <v-col class="my-2" cols="12">
            <v-btn @click="changeTheme" x-large icon>
              <v-icon :color="themeColor" x-large>{{themeOption}}</v-icon>
            </v-btn>
          </v-col>

          <v-col class="my-2" cols="12">
            <v-btn x-large icon>
              <v-icon x-large>mdi-cog</v-icon>
            </v-btn>
          </v-col>

          <v-col class="my-2" cols="12">
            <v-btn @click="logOff" x-large icon>
              <v-icon x-large>mdi-power-standby</v-icon>
            </v-btn>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </v-card>
</template>

<script lang="ts">
import Vue from "vue";
import vuetify from "@/plugins/vuetify";

export default Vue.extend({
  name: "InnerSideBar",
  data: () => ({
    themeOption: "mdi-moon-waxing-crescent",
    themeColor: "",
    highlightedSidebarOption: "Contact",
  }),

  methods: {
    changeTheme: function () {
      vuetify.framework.theme.dark = !vuetify.framework.theme.dark;
      if (vuetify.framework.theme.dark) {
        this.themeColor = "blue";
      } else {
        this.themeColor = "";
      }
    },

    onNotification: function () {
      this.highlightedSidebarOption = "Notif";
    },

    onMessage: function () {
      this.highlightedSidebarOption = "Message";
    },

    onContact: function () {
      this.highlightedSidebarOption = "Contact";
    },

    logOff: function () {
      this.$store.commit("setToken", "");
      this.$router.push("/login");
    },
  },
});
</script>

<style>
</style>