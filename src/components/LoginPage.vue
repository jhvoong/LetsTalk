<template>
  <v-container fill-height>
    <div align="center">
      <v-form>
        <v-row justify="center" align="center">
          <v-col cols="12">
            <v-img src="../assets/unilag.svg" contain height="300" width="300"></v-img>
          </v-col>

          <v-col class="text-center" cols="12">
            <span class="text-center mx-auto" v-if="isAdmin">Welcome Admin</span>
            <span
              class="text-center mx-auto"
              v-else
            >Welcome to the University of Lagos chatting platform</span>
          </v-col>

          <v-col cols="12" sm="6">
            <v-text-field
              v-model="email"
              outlined
              filled
              type="text"
              :rules="emailRules"
              name="email"
              required
              placeholder="Email Address"
            ></v-text-field>

            <v-text-field
              type="password"
              v-model="password"
              outlined
              filled
              required
              name="password"
              placeholder="Password"
            ></v-text-field>

            <v-alert v-if="signinErrorDetails" dense text type="error">{{signinErrorDetails}}</v-alert>
          </v-col>
          <v-col cols="12">
            <v-btn tile>Log in</v-btn>
          </v-col>
        </v-row>
      </v-form>
    </div>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "LoginPage",

  props: {
    isAdmin: Boolean,
  },

  data: () => ({
    signinErrorDetails: "",
    email: "",
    password: "",

    emailRules: [
      (v: string) => !!v || "E-mail is required",
      (v: string) => /.+@.+\..+/.test(v) || "E-mail must be valid",
    ],
  }),

  methods: {},
});
</script>

<style>
</style>