import axios from "axios";
import { Toast } from "vant";
import router from "@/router/index";
import { ObjMap } from "@/plugin/utils/user";
import { Module } from "vuex";
import { RootState } from "./index.d";
import store from ".";

export interface UserState {
  auth: any;
  token: string;
  userCache: Map<number, any>;
}

const state: UserState = {
  auth: null,
  token: "",
  userCache: new Map<number, any>(),
};

const mutations = {
  setAuth: function (state, user) {
    state.auth = user;
  },
  setToken: function (state, token) {
    state.token = token;
  },
  appendUsers: function (state, users) {
    for (const user of users) {
      state.userCache.set(user.id, user);
    }
  },
};

const actions = {
  async getAuth({ state, commit, rootState }) {
    if (state.auth) return;
    const token = localStorage.getItem("token");
    if (token) {
      commit("setToken", token);
      const res = await axios.get(`/api/v1/auth`);
      // 跟后端的初始化配合
      if (res.data.code === 0) commit("setAuth", res.data.details);
    }
  },
  async login({ state, commit, rootState }, params) {
    try {
      const { data } = await axios.post("/api/v1/user/login", params);
      if (data.code && data.code !== 0) Toast.fail(data.message);
      else {
        commit("setAuth", data.details.user);
        commit("setToken", data.details.token);
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
  async signup({ state, commit, rootState }, params) {
    try {
      const { data } = await axios.post("/api/v2/user", params);
      if (data.code && data.code !== 0) Toast.fail(data.message);
      else {
        commit("setAuth", data.details.user);
        commit("setToken", data.details.token);
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
  async appendUsers({ state, commit, rootState }, ids: number[]) {
    const noExistsId: number[] = [];
    ids.forEach((value) => {
      if (!state.userCache.has(value)) noExistsId.push(value);
    });
    if (noExistsId.length > 0) {
      const { data } = await axios.post(`/api/v1/user/baseUserList`, {
        ids: noExistsId,
      });
      if (data.code && data.code !== 0) Toast.fail(data.message);
      else commit("appendUsers", data.details.list);
    }
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
