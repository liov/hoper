import { createSSRApp } from 'vue'
import App from './App.vue'
import * as Pinia from 'pinia'
import i18n from './locale/index'
import { routeInterceptor } from './interceptors'
import { prototypeInterceptor } from '@hopeio/utils/plugin';
import 'virtual:uno.css'
import '@/style/index.scss'
import { store } from './store'

export function createApp() {
  const app = createSSRApp(App)
  app.use(store)
  app.use(i18n)
  app.use(routeInterceptor)
  app.use(prototypeInterceptor)
  return {
    app,
    Pinia, // 此处必须将 Pinia 返回
  }
}
