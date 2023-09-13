import { createRouter, createWebHistory } from "vue-router";
import { getSessionID } from "../utils";

import Login from "../pages/Login.vue";
import User from "../pages/User.vue";

export default createRouter({
  history: createWebHistory(),
  routes: [
    {
      name: "Index",
      path: "/",
      beforeEnter: (_to, _from, next) =>
        next(getSessionID() ? { path: "/_" } : { path: "/login" }),
    },
    {
      name: "Login",
      path: "/login",
      component: Login,
      beforeEnter: (_to, _from, next) =>
        next(getSessionID() ? { path: "/_" } : null),
    },
    {
      name: "User",
      path: "/:username",
      component: User,
      beforeEnter: (_to, _from, next) =>
        next(getSessionID() ? null : { path: "/login" }),
    },
  ],
});
