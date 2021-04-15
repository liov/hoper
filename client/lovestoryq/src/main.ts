import { createApp } from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import store from "./store";
import "@/plugin/axios";
/// <reference path = "./plugin/plugin.d.ts" />;
import HoperPlugin from "@/plugin/plugin";
import * as Vant from "vant";

createApp(App)
  .use(store)
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
  .use(Vant.Toast)
  .use(Vant.Notify)
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
  .use(Vant.Dialog)
  .use(Vant.ActionSheet)
  .use(Vant.ShareSheet)
  .use(Vant.Overlay)
  .use(Vant.Loading)
  .use(Vant.Calendar)
  .use(Vant.DatetimePicker)
  .mount("#app");
