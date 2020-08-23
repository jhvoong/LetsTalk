<template>
  <div
    align="center"
    style="display: grid; align-items: center; height: 100%; justify-content: center;"
  >
    <v-form v-model="detailsValid">
      <v-row class="fill-height" justify="center" align="center">
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

          <v-alert v-if="signinErrorDetails" dense text type="error">{{signinErrorDetails}}</v-alert>
        </v-col>
        <v-col cols="12">
          <v-btn @click="loginUser()" :disabled="detailsDisabled || !detailsValid" tile>Log in</v-btn>
        </v-col>
      </v-row>
    </v-form>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Axios from "axios";
import router from "@/router";
import store from "@/store";

export default Vue.extend({
  name: "Login",

  props: {
    isAdmin: Boolean,
  },

  data: () => ({
    detailsValid: false,
    detailsDisabled: false,

    signinErrorDetails: "",
    email: "",
    password: "",

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

      console.log(URLs[0]);

      Axios.post(URLs[0], "", {
        params: {
          username: this.email,
          password: this.password,
        },

        withCredentials: true,
      })
        .then((Response) => {
          console.log(Response.data);
          if (Response.status == 200) {
            this.$store.commit("setToken", Response.data);
            this.$router.push("/");
          }
        })
        .catch((Error) => {
          console.log(Error);
          this.detailsDisabled = false;

          if (Error.response) {
            this.signinErrorDetails = Error.response.data;
          }
        });
    },
  },
});
</script>

<style>
</style>