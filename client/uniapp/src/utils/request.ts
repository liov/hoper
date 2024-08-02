/* eslint-disable no-param-reassign */
import qs from 'qs'
export type UniRequestOptions = Omit<UniApp.RequestOptions, 'url'> & {
  baseUrl?: string
  query?: Record<string, any>
  /** 出错时是否隐藏错误提示 */
  hideErrorToast?: boolean
}

type UniRequestInterceptor = (options: UniApp.RequestOptions) => UniApp.RequestOptions
type UniResponseInterceptor = (
  response: UniApp.RequestSuccessCallbackResult,
) => UniApp.RequestSuccessCallbackResult
type UniResponseErrorInterceptor = (
  err: UniApp.GeneralCallbackResult,
) => UniApp.GeneralCallbackResult

class UniRequest {
  constructor(defaultConfig?: UniRequestOptions) {
    if (defaultConfig) {
      this.defaults = Object.assign(this.defaults, defaultConfig)
    }
  }

  // 默认的请求配置
  public defaults: UniRequestOptions = {
    baseUrl: '',
    header: {},
    dataType: 'json',
    // #ifndef MP-WEIXIN
    responseType: 'json',
    // #endif
    timeout: 30000,
  }

  // 请求拦截器
  private requestInterceptors = [] as UniRequestInterceptor[]
  // 响应拦截器
  private responseInterceptors = [] as UniResponseInterceptor[]
  // 响应错误拦截器
  private responseErrorInterceptors = [] as UniResponseErrorInterceptor[]
  public interceptors = {
    request: {
      use: (ri: UniRequestInterceptor) => {
        this.requestInterceptors.push(ri)
      },
    },
    response: {
      use: (ri: UniResponseInterceptor, ei: UniResponseErrorInterceptor) => {
        this.responseInterceptors.push(ri)
        this.responseErrorInterceptors.push(ei)
      },
    },
  }

  // 发起请求，默认配置是defaultConfig，也可以传入config参数覆盖掉默认配置中某些属性
  public request<
    T extends string | AnyObject | ArrayBuffer = any,
    D extends string | AnyObject | ArrayBuffer | undefined = any,
  >(
    method: 'GET' | 'POST' | 'PUT' | 'DELETE',
    url: string,
    data?: D,
    config?: UniRequestOptions,
  ): Promise<ResData<T>> {
    return new Promise<ResData<T>>((resolve, reject) => {
      // 接口请求支持通过 query 参数配置 queryString
      if (config.query) {
        const queryStr = qs.stringify(config.query)
        if (url.includes('?')) {
          url += `&${queryStr}`
        } else {
          url += `?${queryStr}`
        }
      }
      url = (config?.baseUrl || this.defaults.baseUrl) + url
      const header = config?.header
        ? Object.assign(this.defaults.header, config.header)
        : this.defaults.header
      const timeout = config?.timeout ? config.timeout : this.defaults.timeout
      let options: UniApp.RequestOptions = {
        ...config,
        method,
        url,
        header,
        timeout,
        data,
        success: (res) => {
          // 执行响应拦截
          for (const ri of this.responseInterceptors) {
            res = ri(res)
          }
          // 状态码 2xx，参考 axios 的设计
          if (
            res.statusCode >= 200 &&
            res.statusCode < 300 &&
            (res.data as ResData<T>).code === 0
          ) {
            // 2.1 提取核心数据 res.data
            resolve(res.data as ResData<T>)
          } else {
            // 其他错误 -> 根据后端错误信息轻提示
            !config.hideErrorToast &&
              uni.showToast({
                icon: 'none',
                title: (res.data as ResData<T>).msg || '请求错误',
              })
            if (res.statusCode === 401) {
              // 401错误  -> 清理用户信息，跳转到登录页
              // userStore.clearUserInfo()
              // uni.navigateTo({ url: '/pages/login/login' })
            }
            reject(res)
          }
        },
        fail: (err) => {
          // 执行响应错误拦截
          for (const ei of this.responseErrorInterceptors) {
            err = ei(err)
          }
          uni.showToast({
            icon: 'none',
            title: '网络错误，换个网络试试',
          })
          reject(err)
        },
      }

      // 执行请求拦截器
      for (const ri of this.requestInterceptors) {
        options = ri(options)
      }
      // 发送请求
      uni.request(options)
    })
  }

  // 发起get请求
  public get<
    T extends string | AnyObject | ArrayBuffer = any,
    D extends string | AnyObject | ArrayBuffer | undefined = any,
  >(url: string, params?: D, config?: UniRequestOptions) {
    return this.request<T, D>('GET', url, params, config)
  }

  // 发起post请求
  public post<
    T extends string | AnyObject | ArrayBuffer = any,
    D extends string | AnyObject | ArrayBuffer | undefined = any,
  >(url: string, data?: D, config?: UniRequestOptions) {
    return this.request<T, D>('POST', url, data, config)
  }
}

const request = new UniRequest({
  baseUrl: '',
  header: {
    'content-type': 'application/json',
    'user-agent': 'uniapp-' + uni.getAppBaseInfo().appName,
  },
  timeout: 10000,
})

export default request
