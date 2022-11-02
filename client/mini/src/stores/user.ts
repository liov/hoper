import axios from "axios";
import { Toast } from "@/plugins/compatible/toast";


import { defineStore } from "pinia";
import Taro from "@tarojs/taro";

export interface UserState {
  auth: any;
  token: string;
  userCache: Map<number, any>;
}

export const state: UserState = {
  auth: null,
  token: "",
  userCache: new Map<number, any>(),
};

const getters = {
  getUser: (state) => (id) => {
    return state.userCache.get(id);
  },
};

const actions = {
  async getAuth() {
    if (state.auth) return;
    const token = localStorage.getItem("token");
    if (token) {
      state.token = token;
      const res = await axios.get(`/api/v1/auth`);
      // 跟后端的初始化配合
      if (res.data.code === 0) state.auth = res.data.details;
    }
  },
  async login(params) {
    try {
      const {
        data: { details },
      } = await axios.post("/api/v1/user/login", params);
      state.auth = details.user;
      state.token = details.token;
      localStorage.setItem("token", details.token);
      axios.defaults.headers["Authorization"] = details.token;
      await Taro.navigateTo({url:"/pages/index/index"});
    } catch (error: any) {
      if (error.response && error.response.status === 401) {
        throw new Error("Bad credentials");
      }
      throw error;
    }
  },
  async signup(params) {
    try {
      const {
        data: { details },
      } = await axios.post("/api/v2/user", params);
      state.auth = details.user;
      state.token = details.token;
      localStorage.setItem("token", details.token);
      axios.defaults.headers["Authorization"] = details.token;
      await Taro.navigateTo({url:"/pages/index/index"});
    } catch (error: any) {
      if (error.response && error.response.status === 401) {
        throw new Error("Bad credentials");
      }
      throw error;
    }
  },
  async appendUsersById(ids: number[]) {
    const noExistsId: number[] = [];
    ids.forEach((value) => {
      if (!state.userCache.has(value)) noExistsId.push(value);
    });
    if (noExistsId.length > 0) {
      const { data } = await axios.post(`/api/v1/user/baseUserList`, {
        ids: noExistsId,
      });
      if (data.code && data.code !== 0) Toast.fail(data.message);
      else this.appendUsers(data.details.list);
    }
  },
  appendUsers(users) {
    for (const user of users) {
      state.userCache.set(user.id, user);
    }
  },
};

export const useUserStore = defineStore({
  id: "user",
  state: () => state,
  getters,
  actions,
});
