import { createRouter, createWebHashHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import { completedAuthenticated } from "@/router/middle";
import { userRoute } from "@/router/user";
import { momentRoute } from "@/router/moment";
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
];

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: routes.concat(userRoute).concat(momentRoute),
});

export default router;
