import {
  createRouter,
  createWebHashHistory,
  NavigationGuard,
  RouteRecordRaw,
} from "vue-router";
import { completedAuthenticated } from "@/router/middle";
import { userRoute } from "@/router/user";
import { momentRoute } from "@/router/moment";

declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
  }
}

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Index",
    component: () => import("../views/moment/Moment.vue"),
  },
  {
    path: "/chat",
    beforeEnter: completedAuthenticated,
    component: () => import("../views/chat/index.vue"),
  },
  {
    path: "/me",
    name: "Home",
    beforeEnter: completedAuthenticated,
    component: () => import("../views/user/Home.vue"),
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes: routes.concat(userRoute).concat(momentRoute),
});

export default router;
