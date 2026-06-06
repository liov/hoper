import { contentRoute } from "@/router/enum";
import router from "@/router/index";
import emitter from "@/plugin/emitter";

import { contentMutations } from "@/store/content";
import { defineAsyncComponent } from "vue";
import {APP_PLATFORM} from "@/plugin/config";

import {
  createWebHashHistory,
  createWebHistory,
  type RouteRecordRaw,
  type RouterHistory
} from "vue-router";


// 动态路由


export const jump = (path: string, type: number, content: any) => {
  const route = `/${contentRoute[type]}/${content.id}`;
  if (path !== route) {
    contentMutations[type](content);
    router.push(route);
  }
  emitter.emit("onComment");
};

export const _import = (path:string) =>
  defineAsyncComponent(() => import(`../${APP_PLATFORM}/views/${path}.vue`));


const IFrame = () => import("@/layout/frame.vue");
// https://cn.vitejs.dev/guide/features.html#glob-import
const modulesRoutes = import.meta.glob("/src/views/**/*.{vue,tsx}");




function addPathMatch() {
  if (!router.hasRoute("pathMatch")) {
    router.addRoute({
      path: "/:pathMatch(.*)",
      name: "pathMatch",
      redirect: "/error/404"
    });
  }
}


/** 获取路由历史模式 https://next.router.vuejs.org/zh/guide/essentials/history-mode.html */
function getHistoryMode(routerHistory): RouterHistory {
  // len为1 代表只有历史模式 为2 代表历史模式中存在base参数 https://next.router.vuejs.org/zh/api/#%E5%8F%82%E6%95%B0-1
  const historyMode = routerHistory.split(",");
  const leftMode = historyMode[0];
  const rightMode = historyMode[1];
  // no param
  if (historyMode.length === 1) {
    if (leftMode === "hash") {
      return createWebHashHistory("");
    } else if (leftMode === "h5") {
      return createWebHistory("");
    }
  } //has param
  else if (historyMode.length === 2) {
    if (leftMode === "hash") {
      return createWebHashHistory(rightMode);
    } else if (leftMode === "h5") {
      return createWebHistory(rightMode);
    }
  }
  return createWebHashHistory("");
}

export {
  addPathMatch,
  getHistoryMode,
};
