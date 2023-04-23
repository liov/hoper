import { createRouter, createWebHashHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import { completedAuthenticated } from "@/router/middle";

import { defineAsyncComponent } from "vue";
import articleRoute from "@/router/article";
import userRoute from "@/router/user";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Index",
    component: () => import("../views/Index.vue"),
  },
];

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: routes.concat(articleRoute).concat(userRoute),
});

export default router;
