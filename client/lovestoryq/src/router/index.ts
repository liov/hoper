import {
  createRouter,
  createWebHashHistory,
  NavigationGuard,
  RouteRecordRaw,
} from "vue-router";
import { completedAuthenticated } from "@/router/middle";
import { userRoute } from "@/router/user";
import { momentRoute } from "@/router/moment";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Index",
    component: () => import("../views/moment/Moment.vue"),
  },
  {
    path: "/about",
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
