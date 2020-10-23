import Vue from 'vue'
import Vuex from 'vuex'
import createPersistedState from "vuex-persistedstate";

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    token: "",
    email: "",
    name: "",
  },
  plugins: [createPersistedState()],
  mutations: {
    setToken(state, token: string) {
      state.token = token
    },

    setEmail(state, email: string) {
      state.email = email
    }
  },
  actions: {
  },
  modules: {
  }
})
