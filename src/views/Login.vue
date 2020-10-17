<template>
  <div
    align="center"
    style="
      display: grid;
      align-items: center;
      height: 100%;
      justify-content: center;
    "
  >
    <v-form v-model="detailsValid">
      <v-row class="fill-height" justify="center" align="center">
        <v-col cols="12">
          <v-img
            src="../assets/unilag.svg"
            contain
            height="300"
            width="300"
          ></v-img>
        </v-col>

        <v-col cols="12">
          <span v-if="isAdmin">Welcome Admin</span>
          <span v-else
            >Welcome to the University of Lagos chatting platform</span
          >
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
            :disabled="detailsDisabled"
          ></v-text-field>

          <v-text-field
            type="password"
            v-model="password"
            outlined
            filled
            required
            :rules="passwordRules"
            name="password"
            placeholder="Password"
            :disabled="detailsDisabled"
          ></v-text-field>

          <v-alert v-if="signinErrorDetails" dense text type="error">{{
            signinErrorDetails
          }}</v-alert>
        </v-col>
        <v-col cols="12">
          <v-btn
            @click="$router.push('/register')"
            color="green"
            class="mr-4"
            tile
            dark
            >Register</v-btn
          >
          <v-btn
            @click="loginUser()"
            :disabled="detailsDisabled || !detailsValid"
            tile
            >Log in</v-btn
          >
        </v-col>
      </v-row>
    </v-form>
  </div>
</template>

<script lang="ts">
"use strict";

import Vue from "vue";
import Axios from "axios";

export default Vue.extend({
  name: "Login",

  props: {
    isAdmin: Boolean,
  },

  data: () => ({
    detailsValid: false,
    detailsDisabled: false,

    signinErrorDetails: "",
    password: "",
    email: "",

    emailRules: [
      (v: string) => !!v || "E-mail is required",
      (v: string) => /.+@.+\..+/.test(v) || "E-mail must be valid",
    ],
    passwordRules: [(v: string) => !!v || "Password is required"],
  }),

  methods: {
    loginUser: function () {
      console.log(this.$store.state.token);
      this.detailsDisabled = true;

      let URLs: [string, string];
      if (this.isAdmin) {
        URLs = [
          process.env.VUE_APP_BACKEND_SERVER + "/admin/login",
          process.env.VUE_APP_BACKEND_SERVER + "/admin",
        ];
      } else {
        URLs = [
          process.env.VUE_APP_BACKEND_SERVER + "/login",
          process.env.VUE_APP_BACKEND_SERVER + "/",
        ];
      }

      Axios.post(URLs[0], "", {
        params: {
          username: this.email,
          password: this.password,
        },
      })
        .then((Response) => {
          if (Response.status == 200) {
            this.$store.commit("setToken", Response.data);
            this.$store.commit("setEmail", this.email);

            this.$router.push("/");
          }
        })
        .catch((Error) => {
          this.detailsDisabled = false;

          if (Error.response) {
            this.signinErrorDetails = Error.response.data;
          }
        });
    },
  },
});
</script>