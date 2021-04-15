import {
  createRouter,
  createWebHashHistory,
  NavigationGuard,
  RouteRecordRaw,
} from "vue-router";
import store from "@/store/index";
import axios from "axios";

//鉴权
const authenticated: NavigationGuard = (_to, _from, next) => {
  if (store.state.auth) next();
  else next({ name: "Login", query: { back: _to.path } });
};

const completedAuthenticated: NavigationGuard = async (_to, _from, next) => {
  if (store.state.auth && (store.state.auth as any).avatarUrl) next();
  else {
    const res = await axios.get(`/api/v1/user/0`);
    store.commit("SET_AUTH", res.data.details.user);
    next({ name: "Login", query: { back: _to.path } });
  }
};

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Index",
    component: () => import("../views/moment/Moment.vue"),
  },
  {
    path: "/about",
    component: () => import("../views/About.vue"),
  },
  {
    path: "/me",
    name: "Home",
    beforeEnter: completedAuthenticated,
    component: () => import("../views/user/Home.vue"),
  },
  {
    path: "/user/edit",
    name: "Edit",
    beforeEnter: completedAuthenticated,
    component: () => import("../views/user/Edit.vue"),
  },
  {
    path: "/user/login",
    name: "Login",
    component: () => import("../views/user/login.vue"),
  },
  {
    path: "/user/active/:id/:secret",
    name: "Active",
    component: () => import("../views/user/active.vue"),
  },
  {
    path: "/moment/add",
    //beforeEnter: authenticated,
    component: () => import("../views/moment/add.vue"),
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
