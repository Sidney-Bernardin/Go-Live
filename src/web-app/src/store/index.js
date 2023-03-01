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
    SET_ERR_MSG: (state, errMsg) => {
      if (typeof errMsg == String) {
        state.errMsg = errMsg
        return
      }

      console.error('Unexpected Error', errMsg)
      state.errMsg = 'Unexpected error, try again later.'
    },
    SET_SELF: (state, self) => (state.self = self),
  },

  actions: {
    setErrMsg: ({ commit }, errMsg) => commit('SET_ERR_MSG', errMsg),
    setSelf: ({ commit }, self) => commit('SET_SELF', self),
  },

  modules: { diagnostics },
})

export default store
