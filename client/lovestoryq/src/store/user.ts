import axios from "axios";
import { Toast } from "vant";
import router from "@/router/index";
import { ObjMap } from "@/plugin/utils/user";
import { Module } from "vuex";
import { RootState } from "./index.d";

export interface UserState {
  auth: any;
  token: string;
  userCache: ObjMap<number, any>;
}

const state: UserState = {
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
  async getAuth({ state, commit, rootState }) {
    if (state.auth) return;
    const token = localStorage.getItem("token");
    if (token) {
      commit("SET_TOKEN", token);
      const res = await axios.get(`/api/v1/auth`);
      // 跟后端的初始化配合
      if (res.data.code === 0) commit("SET_AUTH", res.data.details);
    }
  },
  async login({ state, commit, rootState }, params) {
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
  appendUsers({ state, commit, rootState }, users) {
    commit("APPEND_USERS", users);
  },
};

const getters = {
  getAuth(state, getters, rootState) {
    return state.auth;
  },
  getUser: (state, getters, rootState) => (id) => {
    return state.userCache.get(id);
  },
};

export const user: Module<UserState, RootState> = {
  state,
  mutations,
  actions,
  getters,
};
