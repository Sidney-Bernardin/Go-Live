import RoomsService from '../../services/RoomsService'

export default {
  namespaced: true,

  state: {
    websocket: null,
    chatMessages: [],
  },

  getters: {},

  mutations: {
    CONNECT: (state, roomID) => {
      const ws = RoomsService.diagnostics(roomID)
      ws.onerror = (err) => console.log(err)
      ws.onmessage = (msg) => {
        console.log(msg)
      }

      state.websocket = ws
    },
    DISCONNECT: (state) => state.websocket.close(1000, 'Goodbye, World!'),
    SEND: (state, msg) => state.websocket.send(JSON.stringify(msg)),
  },

  actions: {
    connect: ({ commit }, roomID) => commit('CONNECT', roomID),
    disconnect: ({ commit }) => commit('DISCONNECT'),
    send: ({ commit }, msg) => commit('SEND', msg),
  },
}
