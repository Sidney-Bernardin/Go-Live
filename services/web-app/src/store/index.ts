import { createStore } from "vuex"

export interface State {
    session: any,
    user: any,
}

export const store = createStore({
    state: {
        session: undefined,
        user: undefined,
    },
    mutations: {
        setSession: (state, session) => state.session = session,
        setUser: (state, user) => state.user = user,
    },
    actions: {
        async loadUser({ commit }, userID: string) {
            console.log("HIHI")

            const user = await new Promise<any>((resolve, _reject) => {
                setTimeout(() => resolve({
                    username: userID,
                }), 300)
            })

            commit("setUser", user)
        },
    },
})
