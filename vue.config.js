
module.exports = {
  "transpileDependencies": [
    "vuetify"
  ],
  "devServer": {
    "https": {
      key: process.env.KEY,
      cert: process.env.CERT,
    },
    "proxy": process.env.VUE_APP_BACKEND_SERVER
  }
}