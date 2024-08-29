import { createRouter, createWebHashHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import { completedAuthenticated } from "@/mixin/router/middle";
import { userRoute } from "@/mixin/router/user";
import { momentRoute } from "@/mixin/router/moment";
import { defineAsyncComponent } from "vue";
import {APP_PLATFORM} from "@/mixin/plugin/config";

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
/*  {
    path: "/diary",
    name: "Wasm",
    component: () => import(`../${APP_PLATFORM}/components/wasm/wasm.vue`),
  },*/
];

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: routes.concat(userRoute).concat(momentRoute),
});

export default router;
