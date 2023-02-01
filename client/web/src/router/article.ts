import type { RouteRecordRaw } from "vue-router";

const articleRoute: Array<RouteRecordRaw> = [
  {
    path: "/article",
    //beforeEnter: authenticated,
    component: () => import("../views/article/List.vue"),
  },
  {
    path: "/article/:id",
    //beforeEnter: authenticated,
    component: () => import("../views/article/Info.vue"),
  },
  {
    path: "/article/edit",
    //beforeEnter: authenticated,
    component: () => import("../views/article/Edit.vue"),
  },
];

export default articleRoute;
