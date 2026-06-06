import Cookies from "js-cookie";
import { useUserStoreHook } from "@/store/modules/user";
import type { UserInfo } from "@/model/user";
import { storage } from "./stroge";

export interface DataInfo<T> {
  /** token */
  accessToken: string;
  /** `accessToken`的过期时间（时间戳） */
  expires?: T;
  /** 用于调用刷新accessToken的接口时所需的token */
  refreshToken: string;
  user?: UserInfo;
}

export const userKey = "user-info";
export const TokenKey = "authorized-token";
/**
 * 通过`multiple-tabs`是否在`cookie`中，判断用户是否已经登录系统，
 * 从而支持多标签页打开已经登录的系统后无需再登录。
 * 浏览器完全关闭后`multiple-tabs`将自动从`cookie`中销毁，
 * 再次打开浏览器需要重新登录系统
 * */
export const multipleTabsKey = "multiple-tabs";

/** 获取`token` */
export function getToken(): DataInfo<number> {
  // 此处与`TokenKey`相同，此写法解决初始化时`Cookies`中不存在`TokenKey`报错
  return Cookies.get(TokenKey)
    ? JSON.parse(Cookies.get(TokenKey))
    : storage.getItem(userKey);
}

/**
 * @description 设置`token`以及一些必要信息并采用无感刷新`token`方案
 * 无感刷新：后端返回`accessToken`（访问接口使用的`token`）、`refreshToken`（用于调用刷新`accessToken`的接口时所需的`token`，`refreshToken`的过期时间（比如30天）应大于`accessToken`的过期时间（比如2小时））、`expires`（`accessToken`的过期时间）
 * 将`accessToken`、`expires`、`refreshToken`这三条信息放在key值为authorized-token的cookie里（过期自动销毁）
 * 将`avatar`、`name`、、`roles`、`permissions`、`refreshToken`、`expires`这七条信息放在key值为`user-info`的localStorage里（利用`multipleTabsKey`当浏览器完全关闭后自动销毁）
 */
export function setToken(data: DataInfo<Date>) {
  let expires = 0;
  const { accessToken, refreshToken, user } = data;
  expires = new Date(data.expires).getTime(); // 如果后端直接设置时间戳，将此处代码改为expires = data.expires，然后把上面的DataInfo<Date>改成DataInfo<number>即可
  const cookieString = JSON.stringify({ accessToken, expires, refreshToken });

  expires > 0
    ? Cookies.set(TokenKey, cookieString, {
        expires: (expires - Date.now()) / 86400000
      })
    : Cookies.set(TokenKey, cookieString);

  Cookies.set(multipleTabsKey, "true", {});

  function setUserKey({ id, avatar, name, roles, role, permissions }) {
    const userStore = useUserStoreHook();
    userStore.SET_AVATAR(avatar);
    userStore.SET_Id(id);
    userStore.SET_NAME(name);
    userStore.SET_ROLES(roles);
    userStore.SET_ROLE(role);
    userStore.SET_PERMS(permissions);
    storage.setItem(userKey, {
      accessToken,
      expires,
      avatar,
      name,
      roles,
      role,
      permissions
    });
  }

  if (user.name && user.roles) {
    const { id, avatar, name, roles, role, permissions } = user;
    setUserKey({
      id,
      avatar,
      name,
      roles,
      role,
      permissions
    });
  } else {
    const id = storage.getItem<UserInfo>(userKey)?.id ?? 0;
    const avatar = storage.getItem<UserInfo>(userKey)?.avatar ?? "";
    const name = storage.getItem<UserInfo>(userKey)?.name ?? "";
    const roles = storage.getItem<UserInfo>(userKey)?.roles ?? [];
    const role = storage.getItem<UserInfo>(userKey)?.role ?? 0;
    const permissions = storage.getItem<UserInfo>(userKey)?.permissions ?? [];
    setUserKey({
      id,
      avatar,
      name,
      roles,
      role,
      permissions
    });
  }
}

/** 删除`token`以及key值为`user-info`的localStorage信息 */
export function removeToken() {
  Cookies.remove(TokenKey);
  Cookies.remove(multipleTabsKey);
  storage.removeItem(userKey);
}

/** 格式化token（jwt格式） */
export const formatToken = (token: string): string => {
  return "Bearer " + token;
};

/** 是否有按钮级别的权限（根据登录接口返回的`permissions`字段进行判断）*/
export const hasPerms = (value: string | Array<string>): boolean => {
  if (!value) return false;
  const allPerms = "*:*:*";
  const { permissions } = useUserStoreHook();
  if (!permissions) return false;
  if (permissions.length === 1 && permissions[0] === allPerms) return true;
  const isAuths = typeof value === "string"
    ? permissions.includes(value)
    : Array.isArray(value) && value.every(item => permissions.includes(item));
  return isAuths;
};
