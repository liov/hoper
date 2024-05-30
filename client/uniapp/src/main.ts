import { createSSRApp } from 'vue'
import App from './App.vue'
import store, { init as storeInit } from './store'
import { init as uniHttpInit } from '@/service/index'
import i18n from './locale/index'
import { routeInterceptor, requestInterceptor, prototypeInterceptor } from './interceptors'
import 'virtual:uno.css'
import '@/style/index.scss'

export function createApp() {
  const app = createSSRApp(App)
  app.use(store)
  app.use(i18n)
  app.use(routeInterceptor)
  app.use(requestInterceptor)
  app.use(prototypeInterceptor)
  storeInit()
  uniHttpInit()
  return {
    app,
  }
}
