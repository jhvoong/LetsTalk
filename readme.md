# Lets Talk

A chatting application using Pion webRTC and gorilla websocket for text, video, voice and file transfer. [Test server](https://unilag-letstalk.herokuapp.com)


# Preface

Lets Talk is a web chatting platform proposal for the University of Lagos (unilag.edu.ng). Due to the covid-19 pandemic, neccesity of having an online chatting/learning platform between students and lecturers is vivid.

Lets Talk supports multi-room chats between users, desktop sharing (during class video sessions), file transfer (over websocket and RTC), video and voice calls.


# Privacy Features

For users of minor age on the platform, a machine learning model is to be created to scan minors chats and report duly to the universities authority if any form of sexual harrasment is found. This would not breach privacy as only the model scans texts and block faulting user/room. Machine learning isn't perfect so faulting users have the choice to fair hearing from administrators.

Users also can not register to the platform. Registration is done by the administrator who generates login details for students using their student email address.


# v1.0.0 Milestone

- [x] Multiple room support for users using gorilla websocket

- [x] Mobile UI

- [x] Desktop UI

- [x] Seamless websocket connection

- [x] Resumable file transfer

- [ ] Desktop screen sharing support

- [ ] Voice and Video call support

- [ ] Low bandwidth consumption using selective video call transfer [#31](https://github.com/metaclips/LetsTalk/issues/31)

- [ ] Add logging system.

- [ ] Admin portal


# Dependencies

## Backend

 - [Golang][go]
 - [pion/webrtc][pion]
 - [mongodb][mongo]
 - github.com/julienschmidt/httprouter 
 - github.com/gorilla/securecookie
 - github.com/gorilla/websocket

[go]: golang.org

[mongo]: go.mongodb.org/mongo-driver

[pion]: https://github.com/pion/webrtc

See [go.mod](go.mod) for more information

## Frontend

 - Vue
 - Vuetify


# Installation
```bash
git clone https://github.com/metaclips/LetsTalk.git
cd LetsTalk

go run .
```


Login credentials are generated by administrators on the administrators portal `/admin`. Default email and password for administrator is `admin@email.com` and `admin` which can later be disable or changed.
Users will login through generated email addresses and a password of the first name used to register  `surname` with the first letter titled e.g. Kaduru


# Configuration

- To update when voice call is implemented


# Browser support

- To update when voice call is implemented

# License

[Apache 2.0](LICENSE)
