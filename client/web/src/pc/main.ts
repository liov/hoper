import HoperPlugin from "@pc/plugin/plugin";
import "uno.css";
import Antd from "ant-design-vue";
import "ant-design-vue/es/message/style/css";
import router from "@pc/router";
import { createApp } from "vue";
import App from "@pc/App.vue";
import { createPinia } from "pinia";
import { init as axiosInit } from "@pc/plugin/axios";

import { init as globalInit } from "./store";

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
