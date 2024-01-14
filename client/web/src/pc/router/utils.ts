import { contentRoute } from "@pc/router/enum";
import router from "@pc/router/index";
import emitter from "@pc/plugin/emitter";

import { contentMutations } from "@pc/store/content";
import { defineAsyncComponent } from "vue";

export const jump = (path: string, type: number, content: any) => {
  const route = `/${contentRoute[type]}/${content.id}`;
  if (path !== route) {
    contentMutations[type](content);
    router.push(route);
  }
  emitter.emit("onComment");
};

export const _import = (path) =>
  defineAsyncComponent(() => import(`../views/${path}.vue`));
