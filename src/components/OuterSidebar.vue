<template>
  <v-container fluid>
    <v-row style="height: 25vh; justify-content: center;">
      <v-col cols="12" sm="8">
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
          <v-list-item v-for="n in 100" :key="n" @click="loadChatContent(roomID)">
            <v-list-item-avatar>
              <v-badge bordered bottom color="deep-purple accent-4" dot offset-x="10" offset-y="10">
                <v-avatar size="30">
                  <v-img src="https://cdn.vuetifyjs.com/images/lists/2.jpg"></v-img>
                </v-avatar>
              </v-badge>
            </v-list-item-avatar>

            <v-list-item-content>
              <v-list-item-title>
                <v-row no-gutters>
                  <v-col
                    cols="auto"
                    class="d-inline-block text-truncate"
                    align="center"
                    style="width: 70%"
                    justify="start"
                  >User name should be here contactskkskskskkskskskskksksksksk</v-col>
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

<script>
import Vue from "vue";

export default Vue.extend({
  name: "OuterSidebar",
  data: () => ({
    rooms: {},
  }),
});
</script>

<style>
.button_default {
  font-family: sans-serif;
  font-weight: 600;
  font-size: small;
}
</style>