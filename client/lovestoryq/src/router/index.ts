import { createRouter, createWebHashHistory, RouteRecordRaw } from "vue-router";
import Moment from "../views/Moment.vue";
import store from "@/store/index";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Moment",
    component: Moment
  },
  {
    path: "/about",
    name: "About",
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/About.vue")
  },
  {
    path: "/me",
    name: "Home",
    beforeEnter: (to, from, next) => {
      if (store.state.user != null) next();
      else next({ name: "Login" });
    },
    component: () => import(/* webpackChunkName: "about" */ "../views/Home.vue")
  },
  {
    path: "/login",
    name: "Login",
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/user/login.vue")
  }
];

const router = createRouter({
  history: createWebHashHistory(),
  routes
});

export default router;
