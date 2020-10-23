FROM node:latest

# make the ui folder the current working directory
WORKDIR /ui

# install simple http server for serving static content
RUN npm install -g http-server

# copy all files to directory
COPY . .

RUN npm install
# build app for production with minification
RUN npm run build

EXPOSE ${PORT}

CMD [ "http-server", "dist" ]