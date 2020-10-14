import { createStore } from "vuex";
import axios from 'axios';

export default createStore({
  state: {
    user: null,
    token: ''
  },
  mutations: {
    SET_USER: function (state, user) {
      state.user = user
    },
    SET_TOKEN: function (state, token) {
      state.token = token
    }
  },
  actions: {
    async login({ commit }, params) {
      try {
        const { data } = await axios.post('/user/login', params)
        commit('SET_USER', data)
      } catch (error) {
        if (error.response && error.response.status === 401) {
          throw new Error('Bad credentials')
        }
        throw error
      }
    },
  },
  modules: {}
});
