import axios from 'axios'
import { createStore } from 'vuex'
import diagnostics from './modules/diagnostics'

const store = createStore({
  namespaced: true,

  state: {
    errMsg: '',
    self: null,
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
  },

  actions: {
    handleError: ({ commit }, err) => commit('HANDLE_ERROR', err),
    setSelf: ({ commit }, self) => commit('SET_SELF', self),
  },

  modules: { diagnostics },
})

export default store
