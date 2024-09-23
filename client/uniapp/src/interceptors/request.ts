import uniHttp from '@/utils/request'
import { API_HOST } from '@/env/config'
import { useUserStore } from '@/store/user'

export function init() {
  uniHttp.defaults.baseUrl = API_HOST
  const token = uni.getStorageSync('token')
  const userStore = useUserStore()
  uniHttp.defaults.header.Authorization = token || userStore.token

  // 添加请求拦截器
  uniHttp.interceptors.request.use(function (config) {
    // 在发送请求之前做些什么
    return config
  })

  // 添加响应拦截器
  uniHttp.interceptors.response.use(
    function (response): UniApp.RequestSuccessCallbackResult {
      // 对响应数据做点什么
      if (response.statusCode === 200) {
        const data = response.data as AnyObject
        if (data.code >= 1003 && data.code <= 1005) {
          uni.showToast({ title: '请登录', icon: 'exception' })
          const pages = getCurrentPages()
          uni.navigateTo({
            url: 'pages/user/login?back=' + pages[pages.length - 1].route,
          })
        } else if (data.code !== 0) {
          uni.showToast({ title: data.msg, icon: 'error' })
          return response
        }
      }
      return response
    },
    function (error): UniApp.GeneralCallbackResult {
      // 对响应错误做点什么
      uni.showToast({ title: error.errMsg, icon: 'error' })
      return error
    },
  )
}
