import { createStore } from 'vuex'
import RoomsService from '../services/RoomsService'

export default createStore({
  state: {
    errorMessage: '',
    self: null,
    websocket: null,
    currentRoomID: '',
    chatMessages: [],
  },

  mutations: {
    HANDLE_ERROR: (state, err) => {
      if (typeof err == 'string') {
        state.errorMessage = err
        return
      }

      console.error('Unexpected Error', err)
      state.errorMessage = 'An unexpected error has accured.'
    },
    SET_SELF: (state, self) => (state.self = self),
    SET_CURRENT_ROOM_ID: (state, roomID) => (state.currentRoomID = roomID),

    SET_WEBSOCKET: (state, ws) => {
      if (ws == null) state.websocket?.close(1000)
      state.websocket = ws
    },
    ON_MESSAGE: (state, msg) => {
      const data = JSON.parse(msg.data)

      switch (data.type) {
        case 'CHAT':
          state.chatMessages.push({
            userID: data.user_id,
            username: data.username,
            text: data.text,
          })
          break
      }
    },
    SEND_MESSAGE: (state, msg) => state.websocket.send(JSON.stringify(msg)),
  },

  actions: {
    handleError: ({ commit }, err) => commit('HANDLE_ERROR', err),
    setSelf: ({ commit }, self) => commit('SET_SELF', self),

    joinRoom: ({ commit }, roomID) => {
      const ws = RoomsService.joinRoom(roomID)
      ws.onerror = (err) => commit('HANDLE_ERROR', err)
      ws.onmessage = (msg) => commit('ON_MESSAGE', msg)
      ws.onopen = (_) => commit('SET_CURRENT_ROOM_ID', roomID)
      ws.onclose = (_) => commit('SET_CURRENT_ROOM_ID', '')

      commit('SET_WEBSOCKET', ws)
    },
    leaveRoom: ({ commit }) => commit('SET_WEBSOCKET', null),
    sendMessage: ({ commit }, msg) => commit('SEND_MESSAGE', msg),
  },

  modules: {},
})
