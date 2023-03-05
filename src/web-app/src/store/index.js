import { createStore } from 'vuex'
import RoomsService from '../services/RoomsService'

const store = createStore({
  namespaced: true,

  state: {
    errMsg: '',
    self: null,
    chatMessages: [
      {
        id: '64037bc46117a8fdc471c823',
        username: 'ajsdoij',
        text: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaasjdlk',
      },
      {
        id: '64037bc46117a8fdc471c823',
        username: 'ajsdoij',
        text: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaasjdlk',
      },
      {
        id: '64037bc46117a8fdc471c823',
        username: 'ajsdoij',
        text: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaasjdlk',
      },
      {
        id: '64037bc46117a8fdc471c823',
        username: 'ajsdoij',
        text: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaasjdlk',
      },
      {
        id: '64037bc46117a8fdc471c823',
        username: 'ajsdoij',
        text: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaasjdlk',
      },
      {
        id: '64037bc46117a8fdc471c823',
        username: 'ajsdoij',
        text: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaasjdlk',
      },
      {
        id: '64037bc46117a8fdc471c823',
        username: 'ajsdoij',
        text: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaasjdlk',
      },
      {
        id: '64037bc46117a8fdc471c823',
        username: 'ajsdoij',
        text: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaasjdlk',
      },
      {
        id: '64037bc46117a8fdc471c823',
        username: 'ajsdoij',
        text: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaasjdlk',
      },
      {
        id: '64037bc46117a8fdc471c823',
        username: 'ajsdoij',
        text: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaasjdlk',
      }
    ],
  },

  mutations: {
    HANDLE_ERROR: (state, err) => {
      if (typeof err == 'string') {
        state.errorMessage = err
        return
      }

      console.error('Unexpected Error', err)
      state.errorMessage = 'Unexpected error, try again later.'
    },
    SET_SELF: (state, self) => (state.self = self),
    SET_WEBSOCKET: (state, ws) => (state.websocket = ws),
    ON_MESSAGE: (state, msg) => {
      console.log(msg)
    },
  },

  actions: {
    handleError: ({ commit }, err) => commit('HANDLE_ERROR', err),
    setSelf: ({ commit }, self) => commit('SET_SELF', self),
    joinRoom: ({ commit }, roomID) => {
      const ws = RoomsService.joinRoom(roomID)
      ws.onerror = (err) => commit('HANDLE_ERROR', err)
      ws.onmessage = (msg) => commit('ON_MESSAGE', msg)
      commit('SET_WEBSOCKET', ws)
    },
  },

  modules: {},
})

export default store
