import { createRouter, createWebHistory } from "vue-router";
import { getSessionID } from "../utils";

import Login from "../pages/Login.vue";
import User from "../pages/User.vue";

export default createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      name: "Login",
      component: Login,
      beforeEnter: (_to, _from, next) =>
        next(getSessionID() ? { name: "User" } : null),
    },
    {
      path: "/:username",
      name: "User",
      component: User,
      beforeEnter: (_to, _from, next) =>
        next(getSessionID() ? null : { name: "Login" }),
    },
  ],
});
