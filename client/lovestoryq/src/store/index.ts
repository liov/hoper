import { createStore } from "vuex";
import axios from "axios";
import CookieUtils from "@/plugin/utils/cookie";

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
    async getUser({ commit }, store) {
      if (store.state.user) return;

      const token = CookieUtils.getCookie("token", document.cookie);
      if (token) {
        commit("SET_TOKEN", token);
        await axios.get(`/api/user/get`).then(res => {
          // 跟后端的初始化配合
          if (res.status === 200) commit("SET_USER", res.data);
        });
      }
    },
    async login({ commit }, params) {
      try {
        const { data } = await axios.post("/user/login", params);
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
