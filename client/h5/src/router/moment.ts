import { RouteRecordRaw } from "vue-router";
import { completedAuthenticated } from "@/router/middle";

export const momentRoute: Array<RouteRecordRaw> = [
  {
    path: "/moment/add",
    //beforeEnter: authenticated,
    component: () => import("../views/moment/Add.vue"),
  },
  {
    path: "/moment/:id",
    component: () => import("../views/moment/Detail.vue"),
  },
];
