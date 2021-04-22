import {createStore} from 'vuex'
import axios from "axios";

const state = {
  numbers: [1, 2, 3],
  user: null,
  token: ""
}

const mutations = {
  ADD_NUMBER(state, payload) {
    state.numbers.push(payload)
  },
  setAuth: function (state, user) {
    state.auth = user;
  },
  setToken: function (state, token) {
    state.token = token;
  }
}

const actions = {
  addNumber(context, number) {
    context.commit('ADD_NUMBER', number)
  },
  async getAuth({state, commit}, params) {
    if (state.auth) return;
    const token = localStorage.getItem("token");
    if (token) {
      commit("setToken", token);
      const res = await axios.get(`/api/user/get`);
      // 跟后端的初始化配合
      if (res.data.code === 200) commit("setAuth", res.data.data);
    }
  },
  async login({commit}, params) {
    try {
      const {data} = await axios.post("/api/user/login", params);
      commit("setAuth", data);
    } catch (error) {
      if (error.response && error.response.status === 401) {
        throw new Error("Bad credentials");
      }
      throw error;
    }
  }
}

const getters = {
  getNumbers(state) {
    return state.numbers
  }
}

const store = createStore({
  state,
  mutations,
  actions,
  getters
})

export default store
