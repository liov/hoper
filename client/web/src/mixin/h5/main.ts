import HoperPlugin from "@/mixin/plugin/plugin";
import "uno.css";
import * as Vant from "vant";
import "vant/es/toast/style";
import "vant/es/dialog/style";
import "vant/es/notify/style";
import "vant/es/image-preview/style";
import router from "@/mixin/router/index";
import { createApp } from "vue";
import App from "@/mixin/h5/App.vue";
import { createPinia } from "pinia";
import { init as axiosInit } from "@/mixin/plugin/axios";

import { init as globalInit } from "@/mixin/store/index";

export const app = createApp(App)
  .use(createPinia())
  .use(router)
  .use(Vant.Col)
  .use(Vant.Row)
  .use(Vant.NavBar)
  .use(Vant.Tabbar)
  .use(Vant.TabbarItem)
  .use(Vant.Tab)
  .use(Vant.Tabs)
  .use(Vant.Icon)
  .use(Vant.List)
  .use(Vant.Skeleton)
  .use(Vant.Cell)
  .use(Vant.CellGroup)
  .use(Vant.Form)
  .use(Vant.Field)
  .use(Vant.Button)
  .use(Vant.Uploader)
  .use(Vant.Picker)
  .use(Vant.Popup)
  .use(HoperPlugin)
  .use(Vant.RadioGroup)
  .use(Vant.Radio)
  .use(Vant.Lazyload, {
    lazyComponent: true,
  })
  .use(Vant.Image)
  .use(Vant.PullRefresh)
  .use(Vant.ShareSheet)
  .use(Vant.ActionSheet)
  .use(Vant.Overlay)
  .use(Vant.Loading)
  .use(Vant.Calendar)
  .use(Vant.DatePicker)
  .use(Vant.Checkbox)
  .use(Vant.CheckboxGroup)
  .use(Vant.Popover)
  .use(Vant.ConfigProvider);

app.mount("#app");

axiosInit();
globalInit();
