import Vue from 'vue'
import VueRouter, { RouteConfig } from 'vue-router'

import Login from '../views/Login.vue'
import RegistUser from "../views/Register.vue"
import Home from "../views/Home.vue"

Vue.use(VueRouter)

const routes: Array<RouteConfig> = [
  {
    path: '/login',
    name: 'Student Login',
    component: Login,
    props: { isAdmin: false }
  },

  {
    path: '/admin/login',
    name: 'Admin Login',
    component: Login,
    props: { isAdmin: true }
  },

  {
    path: '/register',
    name: "User Registration",
    component: RegistUser,
  },

  {
    path: "/",
    name: "Homepage",
    component: Home,
  },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
