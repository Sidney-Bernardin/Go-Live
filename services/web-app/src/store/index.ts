import { createStore } from "vuex";

export default createStore({
  state: {
    self: null,
    room: null,
  },
  mutations: {
    SET_SELF: (state, self) => (state.self = self),
    SET_ROOM: (state, room) => (state.room = room),
  },
  actions: {
    setSelf: ({ commit }, self) => commit("SET_SELF", self),
    setRoom: ({ commit }, room) => commit("SET_ROOM", room),
  },
});
