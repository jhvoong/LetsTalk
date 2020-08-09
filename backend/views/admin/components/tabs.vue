(%define "tabs"%)
<template>
  <v-card>
    <v-row>
      <v-col cols="12">
        <v-img class="white--text align-end" width="200" height="200" src="/assets/unilag.svg" />
      </v-col>

      <v-col cols="12">
        <div align="center">
          <v-row>
            <v-spacer></v-spacer>
            <v-card-title>Welcome Admin</v-card-title>
            <v-spacer></v-spacer>
          </v-row>
        </div>
      </v-col>

      <v-col cols="12">
        <div>
          <v-tabs grow v-model="tabs" dark>
            <v-tab>ADD User(s)</v-tab>
            <v-tab-item>
              <v-card flat>(%template "addUsers" .%)</v-card>
            </v-tab-item>

            <v-tab>Block User</v-tab>
            <v-tab-item>
              <v-card flat>
                <v-card flat>(%template "blockUser" .%)</v-card>
              </v-card>
            </v-tab-item>

            <v-tab>Scan Messages</v-tab>
            <v-tab-item>
              <v-card flat>(%template "messageScanner" .%)</v-card>
            </v-tab-item>
          </v-tabs>
        </div>
      </v-col>
    </v-row>
  </v-card>
</template>
(%end%)

(%define "tabsData"%)
    data(){
        return{
            tabs: '',
            addUserOption: 'Add single user',
            addUserOptions: ['Add from file', 'Add single user'],
            date: new Date().toISOString().substr(0, 10),
            dateFormatted: this.formatDate(new Date().toISOString().substr(0, 10)),
            datePickerMenu: false,

            usersClass: 'student',
            usersClasses: ['student','staff'],
            faculty: '',
            faculties: ['Engineering', 'Faculty B', 'Faculty C'],
            staffFaculty: ['Teaching', 'Non-Teaching'],
            email: '',
            emailRules: [
                v => !!v || 'email address is required for sign in',
                v => /.+@.+\..+/.test(v) || 'email address must be valid',
            ],
        }
    },

    methods:{
        parseDate (date) {
            if (!date) return null

            const [month, day, year] = date.split('/')
            return `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}`
        },

        formatDate (date) {
            if (!date) return null

            const [year, month, day] = date.split('-')
            return `${month}/${day}/${year}`
        },
    },
    watch: {
      date (val) {
        this.dateFormatted = this.formatDate(this.date)
      },
    },

(%end%)