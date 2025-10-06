import { createRouter, createWebHistory, type RouteLocationNormalizedGeneric } from "vue-router";

import Home from "../pages/Home.vue";
import Login from "../pages/Login.vue";
import User from "../pages/User.vue";
import { store } from "../store";
import NotFound from "../pages/NotFound.vue";


const notFoundOpts = (to: RouteLocationNormalizedGeneric) => ({
    name: "NotFound",
    params: { pathMatch: to.path.split("/").slice(1) },
    query: to.query,
    hash: to.hash,
})

export const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            name: "Home",
            path: "/",
            // redirect: "",
            component: Home,
        },
        {
            name: "Login",
            path: "/login",
            component: Login,
        },
        {
            name: "User",
            path: "/:userID",
            component: User,
            beforeEnter: async (to, _from, next) => {
                store.dispatch("loadUser", to.params.userID)
                if (!store.state.user)
                    next(notFoundOpts(to))
                else
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
