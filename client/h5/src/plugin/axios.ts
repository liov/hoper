import axios from "axios";
import { Toast } from "vant";
import router from "@/router/index";
import { API_HOST } from "@/plugin/config";
import { useUserStore } from "@/store/user";

export function init() {
  axios.defaults.baseURL = API_HOST;
  const token = localStorage.getItem("token");
  const userStore = useUserStore();
  axios.defaults.headers["Authorization"] = token ? token : userStore.token;

  // 添加请求拦截器
  axios.interceptors.request.use(
    function (config) {
      // 在发送请求之前做些什么
      return config;
    },
    function (error) {
      // 对请求错误做些什么
      return Promise.reject(error);
    }
  );

  // 添加响应拦截器
  axios.interceptors.response.use(
    function (response) {
      // 对响应数据做点什么
      if (response.status == 200) {
        if (response.data.code >= 1003 && response.data.code <= 1005) {
          Toast("请登录");
          router.push({
            name: "Login",
            query: { back: router.currentRoute.value.path },
          });
        } else if (response.data.code !== 0) {
          Toast.fail(response.data.message);
          return Promise.reject({ response: response });
        }
      }
      return Promise.resolve(response);
    },
    function (error) {
      // 对响应错误做点什么
      Toast.fail(error);
      return Promise.reject(error);
    }
  );
}
