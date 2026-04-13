import { defineStore } from "pinia";
import {
  store
} from "../utils";
import { storage } from "@/utils/stroge";
import { getLogin, refreshTokenApi } from "@/api/user";
import router from "@/router";
import { setToken, removeToken, userKey, type DataInfo } from "@/utils/auth";
import type { UserInfo } from "@/model/user";


type State = UserInfo;
export const useUserStore = defineStore("user",{
  state: (): State => ({
    id: storage.getItem<State>(userKey)?.id ?? 0,
    // 姓名
    name: storage.getItem<State>(userKey)?.name ?? "",
    phone: storage.getItem<State>(userKey)?.phone ?? "",
    // 页面级别权限
    roles: storage.getItem<State>(userKey)?.roles ?? [],
    role: storage.getItem<State>(userKey)?.role ?? 1,
    // 按钮级别权限
    permissions: storage.getItem<State>(userKey)?.permissions ?? []
  }),
  actions: {
    SET_Id(id: number) {
      this.id = id;
    },
    /** 存储头像 */
    SET_AVATAR(avatar: string) {
      this.avatar = avatar;
    },
    /** 存储姓名 */
    SET_NAME(name: string) {
      this.name = name;
    },

    /** 存储角色 */
    SET_ROLES(roles: Array<string>) {
      this.roles = roles;
    },

    SET_ROLE(role: number) {
      this.role = role;
    },
    /** 存储按钮级别权限 */
    SET_PERMS(permissions: Array<string>) {
      this.permissions = permissions;
    },
    /** 登入 */
    async loginByUsername(data): Promise<DataInfo<Date>> {
      return new Promise<DataInfo<Date>>((resolve, reject) => {
        getLogin(data)
          .then(res => {
            console.log(res)
            setToken(res);
            resolve(res);
          })
          .catch(error => {
            reject(error);
          });
      });
    },
    /** 前端登出（不调用接口） */
    logOut() {
      this.name = "";
      this.roles = [];
      this.permissions = [];
      removeToken();
      router.push("/login");
    },
    /** 刷新`token` */
    async handRefreshToken(data) {
      return new Promise<DataInfo<Date>>((resolve, reject) => {
        refreshTokenApi(data)
          .then(data => {
            if (data) {
              setToken(data);
              resolve(data);
            }
          })
          .catch(error => {
            reject(error);
          });
      });
    }
  }
});

export function useUserStoreHook() {
  return useUserStore(store); // 有参数,为了SSR
}
