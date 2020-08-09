(%define "addUsers"%)
<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12" sm="4">
        <v-select
          v-model="addUserOption"
          :items="addUserOptions"
          outlined
          hint="Select format to add user(s)"
          persistent-hint
          width="300"
        ></v-select>
      </v-col>

      <v-col cols="12">
        <div v-if="addUserOption==='Add from file'">
          <span>Not yet implemented.</span>
        </div>
        <div v-else>
          <v-row justify="center" style="width: 400px;">
            (%if .UploadSuccess%)
            <v-col cols="12">
              <v-alert type="success">User successfully uploaded</v-alert>
            </v-col>(%end%)
            (%if .Error%)
            <v-col cols="12">
              <v-alert type="success">Could not upload user</v-alert>
            </v-col>(%end%)
            <form action="/admin/upload" method="POST">
              <v-col cols="12">
                <v-text-field
                  outlined
                  v-model="email"
                  :rules="emailRules"
                  required
                  name="email"
                  rounded
                  label="Email address*"
                ></v-text-field>
              </v-col>

              <v-col cols="12">
                <v-text-field outlined name="name" rounded label="Full Name. Surname First.*"></v-text-field>
              </v-col>

              <v-col cols="12">
                <v-menu
                  ref="datePickerMenu"
                  v-model="datePickerMenu"
                  :close-on-content-click="false"
                  transition="scale-transition"
                  offset-y
                  max-width="290px"
                  min-width="290px"
                >
                  <template v-slot:activator="{ on }">
                    <v-text-field
                      v-model="dateFormatted"
                      label="Date of birth*"
                      hint="MM/DD/YYYY format"
                      rounded
                      required
                      outlined
                      name="DOB"
                      persistent-hint
                      @blur="date = parseDate(dateFormatted)"
                      v-on="on"
                    ></v-text-field>
                  </template>
                  <v-date-picker v-model="date" no-title @input="datePickerMenu = false"></v-date-picker>
                </v-menu>
              </v-col>

              <v-col cols="12">
                <v-select
                  v-model="usersClass"
                  :items="usersClasses"
                  outlined
                  rounded
                  required
                  name="usersClass"
                  hint="Select users' class*"
                  persistent-hint
                ></v-select>
              </v-col>

              <v-col v-if="usersClass==='student'" cols="12">
                <v-text-field
                  label="Age"
                  type="number"
                  name="age"
                  required
                  rounded
                  outlined
                  persistent-hint
                ></v-text-field>
              </v-col>

              <v-col v-if="usersClass==='student'" cols="12">
                <v-text-field
                  label="Parent email"
                  name="parentEmail"
                  rounded
                  outlined
                  persistent-hint
                ></v-text-field>
              </v-col>

              <v-col v-if="usersClass==='student'" cols="12">
                <v-text-field
                  label="Parent number"
                  name="parentNumber"
                  rounded
                  outlined
                  persistent-hint
                ></v-text-field>
              </v-col>

              <v-col cols="12">
                <v-select
                  v-if="usersClass==='student'"
                  v-model="faculty"
                  :items="faculties"
                  outlined
                  required
                  name="faculty"
                  rounded
                  hint="Faculty*"
                  persistent-hint
                ></v-select>
                <v-select
                  v-else
                  v-model="faculty"
                  :items="staffFaculty"
                  outlined
                  name="faculty"
                  rounded
                  hint="Faculty*"
                  persistent-hint
                ></v-select>
              </v-col>

              <v-col cols="12">
                <v-btn type="submit">Submit</v-btn>
              </v-col>
            </form>
          </v-row>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>
(%end%)
