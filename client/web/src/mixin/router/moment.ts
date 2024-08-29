import type { RouteRecordRaw } from "vue-router";
import {APP_PLATFORM} from "@/mixin/plugin/config";

export const momentRoute: Array<RouteRecordRaw> = [
  {
    path: "/moment/add",
    //beforeEnter: authenticated,
    component: () => import(`../${APP_PLATFORM}/views/moment/Add.vue`),
  },
  {
    path: "/moment/:id",
    component: () => import(`../${APP_PLATFORM}/views/moment/Detail.vue`),
  },
];
