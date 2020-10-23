# Lets Talk

A chatting application using Pion webRTC and gorilla websocket for text, video, voice and file transfer. [Test server](https://unilag-letstalk.herokuapp.com)


# Preface

Lets Talk is a web chatting platform proposal for students. Due to the covid-19 pandemic, neccesity of having an online chatting/learning platform between students and lecturers is vivid.

Lets Talk supports multi-room chats between users, resumable file transfer (over websocket), video and voice calls support with desktop sharing. The idea is to allow minimal bandwidth consumption during video calls using one video - multiple audio during call session.

Administrator can only register user ***tentative***


# v0.1.0 Milestone

- [x] Multiple room support for users using gorilla websocket

- [x] Mobile UI

- [x] Desktop UI

- [x] Seamless websocket connection

- [x] Resumable file transfer

- [x] Seamless Desktop screen sharing support

- [x] Voice and Video call support

- [x] Low bandwidth consumption using selective video call transfer [#31](https://github.com/metaclips/LetsTalk/issues/31)

- [x] Video session upload to cloud platforms or Database.

- [ ] Admin portal


# Dependencies

## Backend

 - [Golang](golang.org)
 - [pion/webrtc](https://github.com/pion/webrtc)
 - [mongodb](go.mongodb.org/mongo-driver)
 - [httprouter](github.com/julienschmidt/httprouter)
 - [gorilla secure cookies](github.com/gorilla/securecookie)
 - [gorilla websocket](github.com/gorilla/websocket)


See [go.mod](go.mod) for more information

## Frontend

 - [Vue](https://vuejs.org)
 - [Vuetify](https://vuetifyjs.com)


# Installation
Clone repository
```bash
git clone https://github.com/metaclips/LetsTalk.git
```

## Backend

```bash
go mod download
go run .
```

## Frontend

```bash
npm install
npm run serve
```

## Docker Installation
There are individual docker files for [API](api.Dockerfile) and [Frontend](ui.Dockerfile), there's also a docker compose file.

```bash
docker-compose up
```

# Configuration

Do check [config](config.json) for supported configurations.

To generate TLS certificate for test using tls.

` openssl req -nodes -x509 -newkey rsa:4096 -keyout key.pem -subj "/C=US/ST=Oregon/L=Portland/O=Company Name/OU=Org/CN=example.com" -out cert.pem -days 365 `

- The default ICE server in use if none is specified in config is: `stun:stun.l.google.com:19302`

- If no TLS key and cert path is specified, HTTP is used.

- If Class session record is enabled, call sessions are recorded and uploaded to Dropbox server if a dropbox token is provided or file database.


# License

[Apache 2.0](LICENSE)
