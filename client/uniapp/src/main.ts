import { createSSRApp } from 'vue'
import App from './App.vue'
import * as Pinia from 'pinia'
import i18n from './locale/index'
import { routeInterceptor, init as uniHttpInit  } from './interceptors'
import { prototypeInterceptor } from '@hopeio/utils/plugin';
import 'virtual:uno.css'
import '@/style/index.scss'

export function createApp() {
  const app = createSSRApp(App)
  app.use(Pinia.createPinia())
  app.use(i18n)
  app.use(routeInterceptor)
  app.use(prototypeInterceptor)
  uniHttpInit()
  return {
    app,
    Pinia, // 此处必须将 Pinia 返回
  }
}
