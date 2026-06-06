import { httpclient } from "./common";
import type { ListResp } from "@hopeio/utils/types";
import type { UserInfo } from "@/model/user";
import type { DataInfo } from "@/utils/auth";
import  {LoginResp} from "@gen/pb/user/user.service";

/** 登录 */
export const getLogin = (data?: object) => {
  return httpclient.request<LoginResp>("post", "/api/login", {
    data,
    decode: LoginResp.decode,
  });
};

/** 刷新`token` */
export const refreshTokenApi = (data?: object) => {
  return httpclient.request<DataInfo<Date>>("post", "/refresh-token", {
    data
  });
};

/** 账户设置-个人信息 */
export const getMine = (data?: object) => {
  return httpclient.request<UserInfo>("get", "/mine", { data });
};

/** 账户设置-个人安全日志 */
export const getMineLogs = (data?: object) => {
  return httpclient.request<ListResp>("get", "/mine-logs", { data });
};

export const getSimpleUsers = (data?: object) => {
  return httpclient.request<ListResp>("get", "/api/simpleUser", {
    params: data,
    responseType:'json'
  });
};

export const resetPassword = (data?: object) => {
  return httpclient.request<ListResp>("put", "/api/user/resetPassword", {
    data
  });
};
