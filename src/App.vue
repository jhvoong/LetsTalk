<template>
  <div id="app">
    <v-app>
      <v-main>
        <router-view></router-view>
      </v-main>
      <Footer />
    </v-app>
  </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "App",

  data: () => ({
    //
  }),
});
</script>

st
