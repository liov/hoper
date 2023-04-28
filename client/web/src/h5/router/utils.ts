import { contentRoute } from "@h5/router/enum";
import router from "@h5/router/index";
import emitter from "@h5/plugin/emitter";

import { contentMutations } from "@h5/store/content";
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
