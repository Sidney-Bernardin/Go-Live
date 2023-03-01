import { createRouter, createWebHistory } from 'vue-router'

import Home from '/src/pages/Home.vue'
import Login from '/src/pages/Login.vue'
import User from '/src/pages/User.vue'

import { getSessionID } from '../utils'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: Home,
      beforeEnter: (to, from, next) =>
        next(getSessionID() ? null : { name: 'Login' }),
    },
    {
      path: '/login',
      name: 'Login',
      component: Login,
      beforeEnter: (to, from, next) =>
        next(getSessionID() ? { name: 'Home' } : null),
    },
    {
      path: '/:username',
      name: 'User',
      component: User,
      beforeEnter: (to, from, next) =>
        next(getSessionID() ? null : { name: 'Login' }),
    },
  ],
})

export default router
