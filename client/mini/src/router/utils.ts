import { contentRoute } from "@/router/enum";
import router from "@/router/index";
import emitter from "@/plugins/emitter";

import { contentMutations } from "@/stores/content";
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
  defineAsyncComponent(() => import(`../pages/${path}.vue`));
