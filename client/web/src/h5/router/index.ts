import { createRouter, createWebHashHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import { completedAuthenticated } from "@h5/router/middle";
import { userRoute } from "@h5/router/user";
import { momentRoute } from "@h5/router/moment";
import { defineAsyncComponent } from "vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Index",
    component: () => import("../views/moment/Moment.vue"),
  },
  {
    path: "/chat",
    name: "Chat",
    beforeEnter: completedAuthenticated,
    component: defineAsyncComponent(() => import("../views/chat/index.vue")),
    meta: { requiresAuth: true },
  },
  {
    path: "/me",
    name: "Home",
    beforeEnter: completedAuthenticated,
    component: () => import("../views/user/Home.vue"),
  },
  {
    path: "/diary",
    name: "Wasm",
    component: () => import("../components/wasm/wasm.vue"),
  },
];

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: routes.concat(userRoute).concat(momentRoute),
});

export default router;
