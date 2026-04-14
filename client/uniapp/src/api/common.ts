import { HttpClient, RequestOptions } from '@hopeio/utils/uniapp'
import { Ref, UnwrapRef } from 'vue'
import { CommonResp } from '@hopeio/utils/types'


const defaultConfig: RequestOptions = {
  baseURL: import.meta.env.VITE_SERVER_BASEURL,
  // 请求超时时间
  timeout: 10000,
  responseType: "arraybuffer",
  headers: {
    Authorization: uni.getStorageSync('token'),
    Accept: "application/json, text/plain, */*",
    "Content-Type": "application/json",
    "X-Requested-With": "XMLHttpRequest"
  },

};
export const httpclient = new HttpClient(defaultConfig);

// 添加请求拦截器
httpclient.interceptors.request.use(function (config) {
  if (JSON.parse(__VITE_APP_PROXY__)) {
    //
  }
  // 在发送请求之前做些什么
  uni.showLoading({
    title: config.loadingMsg || '加载中'
  })
  return config
})

// 添加响应拦截器
httpclient.interceptors.response.use(
  function (res) {
    // 对响应数据做点什么
    if (res.response.statusCode === 200) {
      const data = res.response.data as AnyObject
      if (data.code >= 1003 && data.code <= 1005) {
        uni.showToast({ title: '请登录', icon: 'exception' })
        const pages = getCurrentPages()
        uni.navigateTo({
          url: '/pages/user/login?back=' + pages[pages.length - 1].route,
        })
      } else if (data.code !== 0) {
        uni.showToast({ title: data.msg, icon: 'error' })
        return res
      }
    }
    return res
  },
  function (error): UniApp.GeneralCallbackResult {
    // 对响应错误做点什么
    uni.showToast({ title: error.errMsg, icon: 'error' })
    return error
  },
)


type IUseRequestOptions<T> = {
  /** 是否立即执行，如果是则在onLoad执行 */
  immediate?: boolean
  /** 初始化数据 */
  initialData?: T
}

/**
 * useRequest是一个定制化的请求钩子，用于处理异步请求和响应。
 * @param func 一个执行异步请求的函数，返回一个包含响应数据的Promise。
 * @param options 包含请求选项的对象 {immediate, initialData}。
 * @param options.immediate 是否立即执行请求，默认为true。
 * @param options.initialData 初始化数据，默认为undefined。
 * @returns 返回一个对象{loading, error, data, run}，包含请求的加载状态、错误信息、响应数据和手动触发请求的函数。
 */
export function useRequest<T>(
  func: () => Promise<T>,
  options: IUseRequestOptions<T> = { immediate: true },
): { loading: Ref<boolean>; error: Ref<boolean>; data: Ref<UnwrapRef<T> | undefined>; run: () => Promise<void> } {
  const loading = ref(false)
  const error = ref(false)
  const data = ref<T | undefined>(options.initialData)
  const run = async () => {
    loading.value = true
    func()
      .then((res) => {
        data.value = res as UnwrapRef<T>
        error.value = false
      })
      .catch((err) => {
        error.value = err
      })
      .finally(() => {
        loading.value = false
      })
  }

  onLoad(() => {
    options.immediate && run()
  })
  return { loading: loading as Ref<boolean>, error: error as Ref<boolean>, data: data as Ref<UnwrapRef<T> | undefined>, run }
}
