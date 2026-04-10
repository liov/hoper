import { httpclient } from '@hopeio/utils/uniapp'
import { API_HOST } from '@/env/config'
import { useUserStore } from '@/store/user'

export function json() {
  httpclient.defaults.baseUrl = API_HOST
  const token = uni.getStorageSync('token')
  const userStore = useUserStore()
  httpclient.defaults.header.Authorization = token || userStore.token

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
}


export function protobuf() {
  httpclient.defaults.baseUrl = API_HOST
  const token = uni.getStorageSync('token')
  const userStore = useUserStore()
  httpclient.defaults.header.Authorization = token || userStore.token

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
}
