import { stringify } from "qs";
import NProgress from "@/utils/progress";
import { getToken, removeToken } from "@/utils/auth";
import { useUserStoreHook } from "@/store/modules/user";

import { message } from "@hopeio/utils/message";

import {
  HttpClient,
  type RequestConfig,
  type Response
} from "@hopeio/utils/httpclient";
import axios, { type CustomParamsSerializer } from "axios";
import { ErrResp } from "@gen/pb/hopeio/response/response";
import router from "@/router";

export const baseUrlApi = (url: string) => `/api/${url}`;
// 相关配置请参考：www.axios-js.com/zh-cn/docs/#axios-request-config-1
const defaultConfig: RequestConfig = {
  baseURL: import.meta.env.VITE_API_HOST,
  // 请求超时时间
  timeout: 10000,
  responseType: "arraybuffer",
  responseEncoding: "binary",
  headers: {
    Accept: "application/json, text/plain, */*",
    "Content-Type": "application/json",
    "X-Requested-With": "XMLHttpRequest"
  },
  // 数组格式参数序列化（https://github.com/axios/axios/issues/5142）
  paramsSerializer: {
    serialize: stringify as unknown as CustomParamsSerializer
  }
};
export const httpclient = new HttpClient(defaultConfig);
/** `token`过期后，暂存待执行的请求 */
let requests = [];

/** 防止重复刷新`token` */
let isRefreshing = false;
const errHeader = true;

/** 重连原始请求 */
function retryOriginalRequest(config: RequestConfig) {
  return new Promise(resolve => {
    requests.push((token: string) => {
      config.headers["Authorization"] = token;
      resolve(config);
    });
  });
}

/** 请求拦截 */

httpclient.interceptors.request.use(
  (config: RequestConfig): any => {
    // 开启进度条动画
    NProgress.start();
    // 优先判断post/get等方法是否传入回调，否则执行初始化设置等回调
    if (typeof config.beforeRequestCallback === "function") {
      config.beforeRequestCallback(config);
      return config;
    }

    /** 请求白名单，放置一些不需要`token`的接口（通过设置请求白名单，防止`token`过期后再请求造成的死循环问题） */
    const whiteList = ["/refresh-token", "/api/login"];
    return whiteList.some(url => config.url.endsWith(url))
      ? config
      : new Promise(resolve => {
          const data = getToken();
          if (data) {
            const now = new Date().getTime();
            const expired = parseInt(data.expires) - now <= 0;
            if (expired) {
              if (!isRefreshing) {
                isRefreshing = true;
                // token过期刷新
                useUserStoreHook()
                  .handRefreshToken({ refreshToken: data.refreshToken })
                  .then(res => {
                    const token = res.accessToken;
                    config.headers["Authorization"] = token;
                    requests.forEach(cb => cb(token));
                    requests = [];
                  })
                  .finally(() => {
                    isRefreshing = false;
                  });
              }
              resolve(retryOriginalRequest(config));
            } else {
              config.headers["Authorization"] = data.accessToken;
              resolve(config);
            }
          } else {
            resolve(config);
          }
        });
  },
  error => {
    message(error.message, { type: "error" });
    return error;
  }
);

/** 响应拦截 */

httpclient.interceptors.response.use(
  // @ts-ignore
  (response: Response) => {
    const $config = response.config;
    // 关闭进度条动画
    NProgress.done();
    // 优先判断post/get等方法是否传入回调，否则执行初始化设置等回调
    if (typeof $config.beforeResponseCallback === "function") {
      $config.beforeResponseCallback(response);
      return response;
    }
    let errresp;
    if (errHeader) {
      errresp = {
        code: response.headers["error-code"],
        msg: response.headers["error-msg"]
      };
    } else {
      const contentType = response.headers["content-type"] as string;
      if (
        contentType?.startsWith("application/protobuf") ||
        contentType?.startsWith("application/x-protobuf")
      ) {
        errresp = ErrResp.decode(new Uint8Array(response.data));
      } else {
        errresp = response.data;
      }
    }

    if (errresp?.code) {
      if (errresp.code >= 1002 && errresp.code <= 1005) {
        router.push("/login");
        removeToken();
      }

      message(errresp.msg, { type: "error" });
      return Promise.reject(errresp);
    }
    if ($config.successMsg) {
      message($config.successMsg, { type: "success" });
    }
    return response;
  },
  (error: any) => {
    error.isCancelRequest = axios.isCancel(error);
    // 关闭进度条动画
    NProgress.done();
    message(error.message, { type: "error" });
    // 所有的响应异常 区分来源为取消请求/非取消请求
    return Promise.reject(error);
  }
);
