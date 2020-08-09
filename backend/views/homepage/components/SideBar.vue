(%define "sidebar" %)
<template>
  <div style="max-width: 300px; height: 90%;">
    <v-row>
      <v-col cols="12">
        <v-img src="./assets/unilag.svg" align="center" contain height="150"></v-img>
      </v-col>

      <v-col cols="12" align="center" justify="center">
        <span>Welcome {{name}}</span>
      </v-col>

      <v-col v-if="$vuetify.breakpoint.mdAndUp" cols="12">
        <v-row justify="center" align="center">
          <v-col md="auto">
            <v-dialog max-width="500px" v-model="createRoomDialog">
              <template v-slot:activator="{ on }">
                <v-btn outlined height="50" width="50" v-on="on">
                  <v-icon>mdi-plus</v-icon>
                </v-btn>
              </template>

              <v-card>
                <v-card-text>
                  <v-container fluid>
                    <v-row>
                      <v-col cols="12">
                        <span>Create New Room</span>
                      </v-col>
                      <v-col cols="12">
                        <v-text-field
                          label="Specify Room Name"
                          @keyup.enter.exact="createRoom"
                          v-model="newRoomName"
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12">
                        <v-spacer></v-spacer>
                        <v-btn color="green darken-1" text @click="createRoom">Create Room</v-btn>
                      </v-col>
                    </v-row>
                  </v-container>
                </v-card-text>
              </v-card>
            </v-dialog>
          </v-col>

          <v-col md="auto">
            <v-badge
              @click.native="openMessageDialog=!openMessageDialog"
              :content="messages.length"
              :value="messages.length"
              color="red"
              overlap
            >
              <v-icon large>mdi-email</v-icon>
            </v-badge>

            <v-dialog scrollable v-model="openMessageDialog" width="600px">
              <v-card>
                <v-card-title>Messages</v-card-title>
                <v-divider></v-divider>
                <v-card-text style="max-height: 500px;">
                  <v-container>
                    <span v-if="messages.length==0">No Message</span>
                    <v-row>
                      <v-col v-for="(message,i) in messages" :key="i" cols="12">
                        <span>{{message.requestingUserName}} ({{message.requestingUserID}}) wants to add you to a room ({{message.roomName}})</span>
                        <v-spacer></v-spacer>
                        <v-btn
                          color="green darken-1"
                          @click="acceptJoinRequest(message.roomID,message.roomName,i)"
                          text
                        >Accept</v-btn>
                        <v-btn color="green darken-1" text>Decline</v-btn>
                      </v-col>
                    </v-row>
                  </v-container>
                </v-card-text>
              </v-card>
            </v-dialog>
          </v-col>
        </v-row>
      </v-col>
    </v-row>

    <v-expansion-panels popout>
      <v-expansion-panel>
        <v-expansion-panel-header>
          <v-row align="center">
            <v-col cols="8" md="8" sm="2">Chats</v-col>
            <v-col md="4" sm="1">
              <v-chip
                color="red"
                x-small
                pill
                text-color="white"
              >{{Object.keys(onreadroommessagecount).length}}</v-chip>
            </v-col>
          </v-row>
        </v-expansion-panel-header>

        <v-expansion-panel-content>
          <v-container style="height: 55vh;" class="overflow-y-auto">
            <v-list tile dense three-line>
              <v-list-item-group color="black">
                <v-list-item v-for="(roomID,i) in rooms" :key="i" @click="loadChatContent(roomID)">
                  <v-list-item-avatar>
                    <v-icon large>mdi-account-circle</v-icon>
                  </v-list-item-avatar>

                  <v-list-item-content>
                    <v-list-item-title>
                      <v-row no-gutters>
                        <v-col
                          cols="auto"
                          class="d-inline-block text-truncate"
                          align="center"
                          style="max-width: 100px;"
                          justify="start"
                        >{{chatspreview[roomID]["roomName"]}}</v-col>
                        <v-col
                          cols="mx-auto"
                          align="center"
                          v-if="onreadroommessagecount[roomID]!==undefined"
                          justify="end"
                        >
                          <v-chip color="red" x-small pill text-color="white">1</v-chip>
                        </v-col>
                      </v-row>
                    </v-list-item-title>

                    <v-list-item-subtitle
                      class="d-inline-block text-truncate"
                    >{{chatspreview[roomID]["message"]}}</v-list-item-subtitle>
                  </v-list-item-content>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-container>
        </v-expansion-panel-content>
      </v-expansion-panel>
    </v-expansion-panels>
  </div>
</template>
(%end%)
