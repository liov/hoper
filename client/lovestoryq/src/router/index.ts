import {
  createRouter,
  createWebHashHistory,
  NavigationGuard,
  RouteRecordRaw
} from "vue-router";
import store from "@/store/index";

//鉴权
const authenticated: NavigationGuard = (_to, _from, next) => {
  console.log(store.state.user);
  if (store.state.user != null) next();
  else next({ name: "Login" });
};

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Index",
    component: () => import("../views/Moment.vue")
  },
  {
    path: "/about",
    component: () => import("../views/About.vue")
  },
  {
    path: "/me",
    name: "Home",
    beforeEnter: authenticated,
    component: () => import("../views/Home.vue")
  },
  {
    path: "/login",
    name: "Login",
    component: () => import("../views/user/login.vue")
  },
  {
    path: "/moment/add",
    //beforeEnter: authenticated,
    component: () => import("../views/moment/add.vue")
  }
];

const router = createRouter({
  history: createWebHashHistory(),
  routes
});

export default router;
