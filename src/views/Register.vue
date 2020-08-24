<template>
  <div
    align="center"
    style="display: grid; align-items: center; height: 100%; justify-content: center;"
  >
    <v-form v-model="dataIsValid">
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
            type="email"
            :rules="emailRules"
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
            :rules="fillInputRules"
            placeholder="Password"
            :disabled="detailsDisabled"
          ></v-text-field>

          <v-text-field
            v-model="name"
            outlined
            filled
            type="text"
            required
            :rules="fillInputRules"
            placeholder="Full Name. Surname first."
            :disabled="detailsDisabled"
          ></v-text-field>

          <v-menu
            v-model="datePickerMenu"
            :close-on-content-click="false"
            transition="scale-transition"
            offset-y
            max-width="290px"
            min-width="290px"
          >
            <template v-slot:activator="{ on, attrs }">
              <v-text-field
                outlined
                filled
                v-model="computedDateFormatted"
                label="Date (read only text field)"
                hint="MM/DD/YYYY format"
                persistent-hint
                readonly
                v-bind="attrs"
                v-on="on"
              ></v-text-field>
            </template>
            <v-date-picker v-model="DOB" no-title @input="datePickerMenu = false"></v-date-picker>
          </v-menu>

          <v-alert
            v-if="registrationErrorDetails"
            dense
            text
            type="error"
          >{{registrationErrorDetails}}</v-alert>
        </v-col>

        <v-col cols="12">
          <v-btn class="mr-4" dark @click="$router.push('/login')" color="green" tile>Log In</v-btn>
          <v-btn @click="registerUser()" :disabled="detailsDisabled || !dataIsValid" tile>Register</v-btn>
        </v-col>
      </v-row>
    </v-form>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Axios from "axios";

export default Vue.extend({
  name: "Register",
  data: () => ({
    detailsDisabled: false,
    dataIsValid: false,
    datePickerMenu: false,

    registrationErrorDetails: "",
    email: "",
    password: "",
    name: "",

    DOB: new Date().toISOString().substr(0, 10),

    emailRules: [
      (v: string) => !!v || "E-mail is required",
      (v: string) => /.+@.+\..+/.test(v) || "E-mail must be valid",
    ],
    fillInputRules: [(v: string) => !!v || "Field is required"],
  }),

  computed: {
    computedDateFormatted(): string {
      return this.formatDate(this.DOB);
    },
  },

  methods: {
    registerUser() {
      this.detailsDisabled = true;

      Axios.post(process.env.VUE_APP_BACKEND_SERVER + "/register", "", {
        params: {
          email: this.email,
          password: this.password,
          name: this.name,
          DOB: this.DOB, // For now, we use the - syntax.
        },
      })
        .then((Response) => {
          if (Response.status == 200) {
            this.$router.push("/login");
          }
        })
        .catch((Error) => {
          this.detailsDisabled = false;

          if (Error.response) {
            this.registrationErrorDetails = Error.response.data;
          }
        });
    },

    formatDate(date: string): string {
      if (!date) return "";

      const [year, month, day] = date.split("-");
      return `${month}/${day}/${year}`;
    },

    parseDate(date: string): string {
      if (!date) return "";

      const [month, day, year] = date.split("/");
      return `${year}-${month.padStart(2, "0")}-${day.padStart(2, "0")}`;
    },
  },
});
</script>

<style>
</style>