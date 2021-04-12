import { createStore } from "vuex";
import axios from "axios";
import { Toast } from "vant";
import router from "@/router/index";
import { ObjMap } from "@/plugin/utils/user";

const state = {
  auth: null,
  token: "",
  userCache: new ObjMap(),
};

const mutations = {
  SET_AUTH: function (state, user) {
    state.auth = user;
  },
  SET_TOKEN: function (state, token) {
    state.token = token;
  },
  APPEND_USERS: function (state, users) {
    state.userCache.append(users);
  },
};

const actions = {
  async getAuth({ state, commit }) {
    if (state.auth) return;
    const token = localStorage.getItem("token");
    if (token) {
      commit("SET_TOKEN", token);
      const res = await axios.get(`/api/v1/auth`);
      // 跟后端的初始化配合
      if (res.data.code === 0) commit("SET_AUTH", res.data.details);
    }
  },
  async login({ commit }, params) {
    try {
      const { data } = await axios.post("/api/v1/user/login", params);
      if (data.code && data.code !== 0) Toast.fail(data.message);
      else {
        commit("SET_AUTH", data.details.user);
        commit("SET_TOKEN", data.details.token);
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
  },
  appendUsers({ commit }, users) {
    commit("APPEND_USERS", users);
  },
};

const getters = {
  getAuth(state) {
    return state.auth;
  },
  getUser: (state) => (id) => {
    return state.userCache.get(id);
  },
};

const store = createStore({
  state,
  mutations,
  actions,
  getters,
});

export default store;
