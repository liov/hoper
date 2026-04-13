import { createRouter, createWebHashHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import { completedAuthenticated } from "@/router/middle";
import { userRoute } from "@/router/user";
import { momentRoute } from "@/router/moment";
import { defineAsyncComponent } from "vue";
import { APP_PLATFORM } from "@/plugin/config";
import { useUserStore } from "@/store/user";

/** 不需要登录可以访问的路由名称 */
const whiteList = ["Login", "Active"];

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Index",
    component: () => import(`../${APP_PLATFORM}/views/index.vue`),
  },
  {
    path: "/chat",
    name: "Chat",
    beforeEnter: completedAuthenticated,
    component: defineAsyncComponent(() => import(`../${APP_PLATFORM}/views/chat/index.vue`)),
    meta: { requiresAuth: true },
  },
  {
    path: "/me",
    name: "Home",
    beforeEnter: completedAuthenticated,
    component: () => import(`../${APP_PLATFORM}/views/user/Home.vue`),
  },
];

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: routes.concat(userRoute).concat(momentRoute),
});

router.beforeEach(async (to, _from, next) => {
  const store = useUserStore();
  if (!store.auth && !store.token) {
    await store.getAuth();
  }
  if (store.auth) {
    // 已登录访问登录页，跳回首页
    if (to.name === "Login") {
      next({ path: "/" });
    } else {
      next();
    }
  } else {
    // 未登录只允许白名单路由
    if (whiteList.includes(to.name as string)) {
      next();
    } else {
      next({ name: "Login", query: { back: to.path } });
    }
  }
});

export default router;
