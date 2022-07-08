import type { RouteRecordRaw } from "vue-router";

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
