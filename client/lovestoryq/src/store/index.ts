import { createStore } from "vuex";
import axios from "axios";
import {Toast} from "vant"
import router from "@/router/index"

const state = {
  user: null,
  token: ""
}

const mutations = {
  SET_USER: function (state, user) {
    state.user = user;
  },
  SET_TOKEN: function (state, token) {
    state.token = token;
  }
}

const actions = {
  async getUser({state, commit}, params) {
    if (state.user) return;
    const token = localStorage.getItem("token");
    if (token) {
      commit("SET_TOKEN", token);
      const res = await axios.get(`/api/v1/auth`);
      // 跟后端的初始化配合
      if (res.data.code === 200) commit("SET_USER", res.data.data);
    }
  },
  async login({commit}, params) {
    try {
      const {data} = await axios.post("/api/v1/user/login", params);
      if (data.code && data.code !== 0) Toast.fail(data.message)
      else {
        store.commit("SET_USER", data.details.user);
        store.commit("SET_TOKEN", data.details.token);
        localStorage.setItem("token", data.details.token);
        axios.defaults.headers["Authorization"] = data.details.token;
        await router.push("/");
      }
    } catch (error) {
      if (error.response && error.response.status === 401) {
        throw new Error("Bad credentials");
      }
      throw error;
    }
  }
}

const getters = {
  getUser(state) {
    return state.user
  }
}

const store = createStore({
  state,
  mutations,
  actions,
  getters
})

export default store
