import { createApp } from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import store from "./store";
import { Col, Row, NavBar, Tabbar, TabbarItem, Tab, Tabs, Icon  } from "vant";

createApp(App)
  .use(store)
  .use(router)
  .use(Col)
  .use(Row)
  .use(NavBar)
  .use(Tabbar)
  .use(TabbarItem)
  .use(Tab)
  .use(Tabs)
  .use(Icon)
  .mount("#app");
