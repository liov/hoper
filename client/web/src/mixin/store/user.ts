import axios, { AxiosError } from "axios";
import { Toast } from "vant";
import router from "@/mixin/router/index";
import { ObjMap } from "@/utils/user";
import { defineStore } from "pinia";

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
  getUser: (state) => (id) => state.userCache.get(id),
};

const actions = {
  async getAuth() {
    if (state.auth) return;
    const token = localStorage.getItem("token");
    if (token) {
      state.token = token;
      const res = await axios.get(`/api/v1/auth`);
      // 跟后端的初始化配合
      if (res.data.code === 0) state.auth = res.data.data;
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
      await router.push("/");
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
      await router.push("/");
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
      if (data.code && data.code !== 0) Toast.fail(data.msg);
      else this.appendUsers(data.data.list);
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
