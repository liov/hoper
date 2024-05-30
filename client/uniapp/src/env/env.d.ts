/// <reference types="vite/client" />
/// <reference types="vite-svg-loader" />

declare module '*.vue' {
  import { DefineComponent } from 'vue'
  // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/ban-types
  const component: DefineComponent<{}, {}, any>
  export default component
}

interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string
  /** 服务端口号 */
  readonly VITE_SERVER_PORT: string
  /** 后台接口地址 */
  readonly VITE_SERVER_BASEURL: string
  /** H5是否需要代理 */
  readonly VITE_APP_PROXY: 'true' | 'false'
  /** H5是否需要代理，需要的话有个前缀 */
  readonly VITE_APP_PROXY_PREFIX: string // 一般是/api
  /** 上传图片地址 */
  readonly VVITE_UPLOAD_BASEURL: string
  /** 是否清除console */
  readonly VITE_DELETE_CONSOLE: 'true' | 'false'
  readonly VITE_SHOW_SOURCEMAP: 'true' | 'false'

  readonly VITE_STATIC_DIR: string
  readonly VITE_API_HOST: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
