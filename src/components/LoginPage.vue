<template>
  <v-container fill-height>
    <div align="center">
      <v-form v-model="detailsValid">
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
              @input="isEmailInvalid"
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
              @input="isPasswordInvalid"
              name="password"
              placeholder="Password"
              :disabled="detailsDisabled"
            ></v-text-field>

            <v-alert v-if="signinErrorDetails" dense text type="error">{{signinErrorDetails}}</v-alert>
          </v-col>
          <v-col cols="12">
            <v-btn
              type="Submit"
              @click="loginUser()"
              :disabled="detailsDisabled || !detailsValid"
              tile
            >Log in</v-btn>
          </v-col>
        </v-row>
      </v-form>
    </div>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Axios from "axios";

export default Vue.extend({
  name: "LoginPage",

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
      this.detailsDisabled = true;

      let URLs: [string, string];
      if (this.isAdmin) {
        URLs = ["/admin/login", "/admin"];
      } else {
        URLs = ["/login", "/"];
      }

      Axios.post(
        URLs[0],
        JSON.stringify({
          username: this.email,
          password: this.password,
        }),
        { withCredentials: true }
      )
        .then((Response) => {
          console.log(Response);
          if (Response.status == 200) {
            this.$router.push(URLs[1]);
          }
        })
        .catch((Error) => {
          console.log(Error);
          this.detailsDisabled = false;
          this.signinErrorDetails = Error.message;
        });
    },
  },
});
</script>

<style>
</style>