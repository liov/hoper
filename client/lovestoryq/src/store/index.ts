import { createStore } from "vuex";
import axios from "axios";

export default createStore({
  state: {
    user: null,
    token: ""
  },
  mutations: {
    SET_USER: function(state, user) {
      state.user = user;
    },
    SET_TOKEN: function(state, token) {
      state.token = token;
    }
  },
  actions: {
    async getUser({ state, commit }, params) {
      if (state.user) return;
      const token = localStorage.getItem("token");
      if (token) {
        commit("SET_TOKEN", token);
        const res = await axios.get(`/api/user/get`);
        // 跟后端的初始化配合
        if (res.data.code === 200) commit("SET_USER", res.data.data);
      }
    },
    async login({ commit }, params) {
      try {
        const { data } = await axios.post("/api/user/login", params);
        commit("SET_USER", data);
      } catch (error) {
        if (error.response && error.response.status === 401) {
          throw new Error("Bad credentials");
        }
        throw error;
      }
    }
  },
  modules: {}
});
