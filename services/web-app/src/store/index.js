import { createStore } from "vuex";

export default createStore({
  state: {
    self: null,
  },
  mutations: {
    SET_SELF: (state, self) => (state.self = self),
  },
  actions: {
    setSelf: ({ commit }, self) => commit("SET_SELF", self),
  },
});
