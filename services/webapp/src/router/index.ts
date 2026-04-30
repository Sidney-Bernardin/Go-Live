import { createRouter, createWebHistory } from "vue-router";

import Login from "../pages/Login.vue";
import User from "../pages/User.vue";
import { store } from "../store";
import NotFound from "../pages/NotFound.vue";


export const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            name: "Login",
            path: "/login",
            component: Login,
            beforeEnter: async (_to, _from, next) => {
                if (store.state.self)
                    next({ name: "Home" })
                else next()
            },
        },
        {
            name: "User",
            path: "/:userID",
            component: User,
            beforeEnter: async (_to, _from, next) => {
                if (!store.state.self)
                    await store.dispatch("getSelf")
                next()
            },
        },
        {
            name: 'NotFound',
            path: '/:pathMatch(.*)*',
            component: NotFound,
        },
    ]
})
