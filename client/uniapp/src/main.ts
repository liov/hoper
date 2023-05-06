import { createSSRApp } from "vue";
import App from "@/App.vue";
import {init as storeInit} from "@/store/index";
import {init as uniHttpInit} from "@/service/index";
import {createPinia} from "pinia";
export function createApp() {
  const app = createSSRApp(App).use(createPinia());
  storeInit();
  uniHttpInit();
  return {
    app,
  };
}

