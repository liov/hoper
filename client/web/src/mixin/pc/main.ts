import HoperPlugin from "@/mixin/plugin/plugin";
import "uno.css";
import Antd from "ant-design-vue";
import router from "@/mixin/router";
import { createApp } from "vue";
import App from "@/mixin/pc/App.vue";
import { createPinia } from "pinia";
import { init as axiosInit } from "@/mixin/plugin/axios";

import { init as globalInit } from "@/mixin/store";

import { ConfigProvider } from "ant-design-vue";

ConfigProvider.config({
  theme: {
    primaryColor: "#25b864",
  },
});
export const app = createApp(App)
  .use(createPinia())
  .use(router)
  .use(Antd)
  .use(HoperPlugin);

app.mount("#app");
//app.config.globalProperties.$message = message;
axiosInit();
globalInit();
