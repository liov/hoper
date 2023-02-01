import type { RouteRecordRaw } from "vue-router";

const userRoute: Array<RouteRecordRaw> = [
  {
    path: "/user/login",
    //beforeEnter: authenticated,
    component: () => import("../views/user/Login.vue"),
  },
  {
    path: "/user/:id",
    //beforeEnter: authenticated,
    component: () => import("../views/article/Info.vue"),
  },
];

export default userRoute;
