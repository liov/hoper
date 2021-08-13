import { RouteRecordRaw } from "vue-router";
import { completedAuthenticated } from "@/router/middle";

export const userRoute: Array<RouteRecordRaw> = [
  {
    path: "/user/edit",
    name: "Edit",
    beforeEnter: completedAuthenticated,
    component: () => import("../views/user/Edit.vue"),
  },
  {
    path: "/user/login",
    name: "Login",
    component: () => import("../views/user/Login.vue"),
  },
  {
    path: "/user/active/:id/:secret",
    name: "Active",
    component: () => import("../views/user/Active.vue"),
  },
];
