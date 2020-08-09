(%define "chattingComponent"%)
<template>
  <div style="height: 100%;" v-if="removewelcomepage===true">
    <v-row no-gutters style="height: 100%;">
      <v-col cols="12">
        <v-container fluid>
          <v-row no-gutters>
            <v-col cols="mx-auto">
              <span>
                <b>{{currentviewedroomname}}</b>
              </span>
            </v-col>
            <v-col cols="auto">
              <v-dialog scrollable v-model="addUserDialog" width="600px">
                <template v-slot:activator="{ on }">
                  <v-btn fab depressed v-on="on">
                    <v-icon>mdi-account-multiple-plus-outline</v-icon>
                  </v-btn>
                </template>

                <v-card>
                  <v-card-title>Add Users</v-card-title>
                  <v-divider></v-divider>
                  <v-card-text style="height: 500px;">
                    <v-row>
                      <v-col cols="12">
                        <v-text-field
                          placeholder="Search Users"
                          rounded
                          v-model="searchText"
                          outlined
                          append-icon="mdi-magnify"
                          @keyup.enter.exact="getUsers"
                          @click:append="getUsers"
                        ></v-text-field>
                      </v-col>

                      <v-col cols="12">
                        <v-checkbox
                          v-for="(user,i) in usersfound"
                          :key="i"
                          :label="user"
                          :value="user"
                          v-model="addedUsers"
                        ></v-checkbox>
                      </v-col>
                    </v-row>
                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="green darken-1" text @click="requestUsersToJoinRoom">Add User(s)</v-btn>
                    <v-btn color="green darken-1" text @click="closeSearchDialog">Close</v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>

              <v-dialog scrollable v-model="showRoomUsersDialog" width="600px">
                <template v-slot:activator="{ on }">
                  <v-btn fab depressed v-on="on">
                    <v-icon>mdi-information</v-icon>
                  </v-btn>
                </template>

                <v-card>
                  <v-card-title>Users</v-card-title>
                  <v-divider></v-divider>
                  <v-card-text style="height: 500px;">
                    <v-row>
                      <v-col cols="12"></v-col>
                      <v-col v-for="(value,key) in onlineusers" :key="key" cols="12">
                        <v-badge inline dot :color="value ? 'green' : 'red'"></v-badge>
                        <span class="mx-4">{{key}}</span>
                      </v-col>
                    </v-row>
                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="red darken-1" text @click="exitRoom">Exit Room</v-btn>
                    <v-btn color="green darken-1" text @click="showRoomUsersDialog=false">Close</v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>

              <v-btn @click="startSession()" fab depressed>
                <v-icon>mdi-phone</v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-container>
        <v-divider></v-divider>
      </v-col>

      <v-col cols="12">
        <v-dialog max-width="300px" persistent v-model="currentchatcontentsloaded">
          <v-card>
            <v-card-text>
              <div align="center" justify="center">
                <v-row>
                  <v-col cols="12">
                    <v-progress-circular indeterminate color="green"></v-progress-circular>
                  </v-col>
                  <v-col cols="12">
                    <span>Fetching Content</span>
                  </v-col>
                </v-row>
              </div>
            </v-card-text>
          </v-card>
        </v-dialog>

        <v-container
          id="chatcontent"
          class="overflow-y-auto scroll-behavior-smooth"
          style="height: 78vh;"
          fluid
        >
          <v-row>
            <v-col cols="12" v-for="(chat,i) in currentchatcontent" :key="i">
              <div align="center" justify="center" v-if="chat['type']==='info'">
                <v-card tile class="justify-center" outlined>
                  <v-card-text>{{chat['message']}}</v-card-text>
                </v-card>
              </div>

              <div align="center" justify="center" v-else-if="chat['type']==='classSession'">
                <v-card tile class="justify-center" outlined>
                  <v-card-text>
                    {{chat['name']}} ({{chat['userID']}}) wants you to join a class session
                    <v-btn text color="green" @click="joinSession(chat)">Click Here To Join</v-btn>
                  </v-card-text>
                </v-card>
              </div>

              <div align="center" justify="center" v-else-if="chat['type']==='classSessionLink'">
                <v-card tile class="justify-center" outlined>
                  <v-card-text>
                    Recorded session by {{chat['name']}} ({{chat['userID']}}) available.
                    <v-btn
                      text
                      color="green"
                      @click="downloadSession(chat.message)"
                    >Click Here To Download</v-btn>
                  </v-card-text>
                </v-card>
              </div>

              <div
                :align="chat.userID==='(%.Email%)' ? 'right' : 'left'"
                justify="center"
                v-else-if="chat['type']==='file'"
              >
                <v-card shaped style="max-width: 70vw;" class="d-inline-block mx-auto">
                  <v-card-title class="text--secondary">
                    <v-row align="center">
                      <v-col :cols="chat.userID==='(%.Email%)' ? 'auto' : 'mx-auto'">
                        <h5>{{chat['name']}}</h5>
                      </v-col>

                      <v-col :cols="chat.userID==='(%.Email%)' ? 'mx-auto' : 'auto'">
                        <v-menu absolute bottom left>
                          <template v-slot:activator="{ on }">
                            <v-btn icon v-on="on">
                              <v-icon>mdi-chevron-down</v-icon>
                            </v-btn>
                          </template>

                          <v-list>
                            <v-list-item v-for="i in 5" :key="i">
                              <v-list-item-title>{{i}}</v-list-item-title>
                            </v-list-item>
                          </v-list>
                        </v-menu>
                      </v-col>
                    </v-row>
                  </v-card-title>

                  <v-card-text>
                    <v-row align="center">
                      <v-col cols="mx-auto">{{chat.message}} ({{chat.fileSize}})</v-col>
                      <v-col cols="auto">
                        <v-progress-circular
                          :rotate="360"
                          :size="50"
                          :width="5"
                          :value="downloadinfo[chat.message].progress"
                          color="teal"
                          v-if="downloadinfo[chat.message]!==undefined&&downloadinfo[chat.message].downloading===true"
                        >
                          <v-btn @click="stopDownload(chat.message)" icon>
                            <v-icon>mdi-close</v-icon>
                          </v-btn>
                        </v-progress-circular>

                        <v-btn @click="startDownload(chat.message,chat.fileHash)" v-else icon>
                          <v-icon>mdi-cloud-download</v-icon>
                        </v-btn>
                      </v-col>
                    </v-row>
                  </v-card-text>
                  <v-card-subtitle>{{chat['time']}}</v-card-subtitle>
                </v-card>
              </div>

              <div align="right" v-else-if="chat['userID']==='(%.Email%)'">
                <v-card shaped style="max-width: 70vw;" class="d-inline-block mx-auto">
                  <v-card-title class="text--secondary">
                    <v-row align="center">
                      <v-col cols="auto">
                        <h5>{{chat['name']}}</h5>
                      </v-col>
                      <v-col cols="mx-auto">
                        <v-menu absolute bottom left>
                          <template v-slot:activator="{ on }">
                            <v-btn icon v-on="on">
                              <v-icon>mdi-chevron-down</v-icon>
                            </v-btn>
                          </template>

                          <v-list>
                            <v-list-item v-for="i in 5" :key="i">
                              <v-list-item-title>{{i}}</v-list-item-title>
                            </v-list-item>
                          </v-list>
                        </v-menu>
                      </v-col>
                    </v-row>
                  </v-card-title>

                  <v-card-text v-if="chat['type']==='upload'">
                    <v-row align="center">
                      <v-col cols="mx-auto">{{chat.message}} ({{chat.fileSize}})</v-col>
                      <v-col cols="auto">
                        <v-btn
                          @click="startDownload(chat.message,chat.fileHash)"
                          depressed
                          text
                          v-if="downloadinfo[chat.message].completed"
                        >
                          <v-icon>mdi-cloud-download</v-icon>
                        </v-btn>

                        <v-progress-circular
                          v-else
                          :rotate="360"
                          :size="50"
                          :width="5"
                          :value="downloadinfo[chat.message].progress"
                          color="teal"
                        >
                          <v-btn
                            @click="stopUpload(chat.message)"
                            v-if="downloadinfo[chat.message].downloading"
                            depressed
                            icon
                          >
                            <v-icon>mdi-close</v-icon>
                          </v-btn>

                          <v-btn icon v-else @click="startUpload(chat.message)">
                            <v-icon>mdi-upload</v-icon>
                          </v-btn>
                        </v-progress-circular>
                      </v-col>
                    </v-row>
                  </v-card-text>

                  <v-card-text v-else align="start">{{chat['message']}}</v-card-text>
                  <v-card-subtitle>{{chat['time']}}</v-card-subtitle>
                </v-card>
              </div>

              <div v-else align="left">
                <v-card shaped style="max-width: 70vw;" class="d-inline-block mx-auto">
                  <v-card-title class="text--secondary">
                    <v-row align="center">
                      <v-col cols="mx-auto">
                        <h5>{{chat['name']}}</h5>
                      </v-col>
                      <v-col cols="auto">
                        <v-menu absolute bottom left>
                          <template v-slot:activator="{ on }">
                            <v-btn icon v-on="on">
                              <v-icon>mdi-chevron-down</v-icon>
                            </v-btn>
                          </template>

                          <v-list>
                            <v-list-item v-for="i in 5" :key="i">
                              <v-list-item-title>{{i}}</v-list-item-title>
                            </v-list-item>
                          </v-list>
                        </v-menu>
                      </v-col>
                    </v-row>
                  </v-card-title>
                  <v-card-text align="start">{{chat['message']}}</v-card-text>
                  <v-card-subtitle align="end">{{chat['time']}}</v-card-subtitle>
                </v-card>
              </div>
            </v-col>
          </v-row>
        </v-container>
      </v-col>

      <v-col cols="12">
        <v-container fluid>
          <v-textarea
            v-model="messageContent"
            prepend-inner-icon="mdi-emoticon"
            prepend-icon="mdi-paperclip"
            append-outer-icon="mdi-send"
            @click:prepend="openFileDialog"
            @click:append-outer="sendMessage"
            @keyup.enter.exact="sendMessage"
            solo
            hide-details="auto"
            dense
            auto-grow
            rows="1"
            rounded
            clearable
            :readonly="disableTextField"
          ></v-textarea>
        </v-container>
      </v-col>
    </v-row>

    <input type="file" id="myFileInput" style="display:none" @change="onFileUpdate" />
  </div>

  <div align="center" v-else>
    <v-row>
      <v-col cols="12">
        <span>We should add something nice here</span>
      </v-col>
      <v-col cols="12">
        <h4>Coming Soon</h4>
      </v-col>
    </v-row>
  </div>
</template>
(%end%)