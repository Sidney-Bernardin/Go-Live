import { createStore } from "vuex"

export interface State {
    session: any,
    user: any,
}

export const store = createStore({
    state: {
        session: undefined,
        self: undefined,
    },
    mutations: {
        setSession: (state, session) => state.session = session,
        setSelf: (state, self) => state.self = self,
    },
    actions: {
        async getSelf({ commit }) {
            var res = await fetch("http://localhost:8000/users/users/self", { credentials: 'include' })
            var json = await res.json()
            if (!res.ok) {
                if (res.status !== 401) {
                    console.error(json)
                    alert("Something went wrong!")
                }
                return
            }

            commit("setSelf", json)
        },
    },
})
