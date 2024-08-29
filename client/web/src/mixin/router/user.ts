import type { RouteRecordRaw } from "vue-router";
import { completedAuthenticated } from "@/mixin/router/middle";
import {APP_PLATFORM} from "@/mixin/plugin/config";

export const userRoute: Array<RouteRecordRaw> = [
  {
    path: "/user/edit",
    name: "Edit",
    beforeEnter: completedAuthenticated,
    component: () => import(`../${APP_PLATFORM}/views/user/Edit.vue`),
  },
  {
    path: "/user/login",
    name: "Login",
    component: () => import(`../${APP_PLATFORM}/views/user/Login.vue`),
  },
  {
    path: "/user/active/:id/:secret",
    name: "Active",
    component: () => import(`../${APP_PLATFORM}/views/user/Active.vue`),
  },
];
