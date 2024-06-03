export async function request<
  D,
  T extends UniNamespace.RequestOptions = UniNamespace.RequestOptions,
>(options: T): Promise<D> {
  return new Promise(function (resolve, reject) {
    uni.request({
      ...options,
      success: (res) => {
        resolve(res.data as D)
      },
      fail: (e) => {
        reject(e)
      },
    })
  })
}

export async function get<D, T extends UniNamespace.RequestOptions = UniNamespace.RequestOptions>(
  options: T,
): Promise<D> {
  return request({
    ...options,
    method: 'GET',
  })
}

export async function post<D, T extends UniNamespace.RequestOptions = UniNamespace.RequestOptions>(
  options: T,
): Promise<D> {
  return request({
    ...options,
    method: 'POST',
  })
}

interface UniResponse<T extends string | AnyObject | ArrayBuffer>
  extends UniApp.RequestSuccessCallbackResult {
  data: T
}
interface UniRequestConfig {
  baseUrl?: string
  /**
   * 设置请求的 header，header 中不能设置 Referer。
   */
  header?: any

  /**
   * 超时时间
   */
  timeout?: number
  /**
   * 如果设为json，会尝试对返回的数据做一次 JSON.parse
   */
  dataType?: string
  /**
   * 设置响应的数据类型。合法值：text、arraybuffer
   */
  responseType?: string
  /**
   * 验证 ssl 证书
   */
  sslVerify?: boolean
  /**
   * 跨域请求时是否携带凭证
   */
  withCredentials?: boolean
  /**
   * DNS解析时优先使用 ipv4
   */
  firstIpv4?: boolean
  /**
   * 开启 http2
   */
  enableHttp2?: boolean
  /**
   * 开启 quic
   */
  enableQuic?: boolean
  /**
   * 开启 cache
   */
  enableCache?: boolean
  /**
   * 是否开启 HttpDNS 服务。如开启，需要同时填入 httpDNSServiceId 。 HttpDNS 用法详见 [移动解析HttpDNS](https://developers.weixin.qq.com/miniprogram/dev/framework/ability/HTTPDNS.html)
   */
  enableHttpDNS?: boolean
  /**
   * HttpDNS 服务商 Id。 HttpDNS 用法详见 [移动解析HttpDNS](https://developers.weixin.qq.com/miniprogram/dev/framework/ability/HTTPDNS.html)
   */
  httpDNSServiceId?: boolean
  /**
   * 开启 transfer-encoding chunked
   */
  enableChunked?: boolean
  /**
   * wifi下使用移动网络发送请求
   */
  forceCellularNetwork?: boolean
  /**
   * 默认 false，开启后可在headers中编辑cookie（支付宝小程序10.2.33版本开始支持）
   */
  enableCookie?: boolean
  /**
   * 是否开启云加速（详见[云加速服务](https://smartprogram.baidu.com/docs/develop/extended/component-codeless/cloud-speed/introduction/)）
   */
  cloudCache?: object | boolean
  /**
   * 控制当前请求是否延时至首屏内容渲染后发送
   */
  defer?: boolean
  success?: (result: RequestSuccessCallbackResult) => void
  /**
   * 失败的回调函数
   */
  fail?: (result: UniApp.GeneralCallbackResult) => void
  /**
   * 结束的回调函数（调用成功、失败都会执行）
   */
  complete?: (result: UniApp.GeneralCallbackResult) => void
}

interface RequestSuccessCallbackResult {
  /**
   * 开发者服务器返回的数据
   */
  data: string | AnyObject | ArrayBuffer
  /**
   * 开发者服务器返回的 HTTP 状态码
   */
  statusCode: number
  /**
   * 开发者服务器返回的 HTTP Response Header
   */
  header: any
  /**
   * 开发者服务器返回的 cookies，格式为字符串数组
   */
  cookies: string[]
}
type UniRequestInterceptor = (options: UniApp.RequestOptions) => UniApp.RequestOptions
type UniResponseInterceptor = (
  response: UniApp.RequestSuccessCallbackResult,
) => UniApp.RequestSuccessCallbackResult
type UniResponseErrorInterceptor = (
  err: UniApp.GeneralCallbackResult,
) => UniApp.GeneralCallbackResult

class UniRequest {
  constructor(defaultConfig?: UniRequestConfig) {
    if (defaultConfig) {
      this.defaults = Object.assign(this.defaults, defaultConfig)
    }
  }

  // 默认的请求配置
  public defaults: UniRequestConfig = {
    baseUrl: '',
    header: {},
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
    config?: UniRequestConfig,
  ): Promise<UniResponse<T>> {
    return new Promise((resolve, reject) => {
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
          resolve(res as UniResponse<T>)
        },
        fail: (err) => {
          // 执行响应错误拦截
          for (const ei of this.responseErrorInterceptors) {
            err = ei(err)
          }
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
  >(url: string, params?: D, config?: UniRequestConfig) {
    return this.request<T, D>('GET', url, params, config)
  }

  // 发起post请求
  public post<
    T extends string | AnyObject | ArrayBuffer = any,
    D extends string | AnyObject | ArrayBuffer | undefined = any,
  >(url: string, data?: D, config?: UniRequestConfig) {
    return this.request<T, D>('POST', url, data, config)
  }
}

const uniHttp = new UniRequest({
  baseUrl: '',
  header: {
    'content-type': 'application/json',
    'user-agent': 'uniapp-' + uni.getAppBaseInfo().appName,
  },
  timeout: 10000,
})

export default uniHttp
