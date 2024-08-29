import type { RouteRecordRaw } from "vue-router";
import {APP_PLATFORM} from "@/mixin/plugin/config";

const articleRoute: Array<RouteRecordRaw> = [
  {
    path: "/article",
    //beforeEnter: authenticated,
    component: () => import(`../${APP_PLATFORM}/views/article/List.vue`),
  },
  {
    path: "/article/:id",
    //beforeEnter: authenticated,
    component: () => import(`../${APP_PLATFORM}/views/article/Info.vue`),
  },
  {
    path: "/article/edit",
    //beforeEnter: authenticated,
    component: () => import(`../${APP_PLATFORM}/views/article/Edit.vue`),
  },
];

export default articleRoute;
