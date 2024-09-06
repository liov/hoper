import { contentRoute } from "@/mixin/router/enum";
import router from "@/mixin/router/index";
import emitter from "@/mixin/plugin/emitter";

import { contentMutations } from "@/mixin/store/content";
import { defineAsyncComponent } from "vue";
import {APP_PLATFORM} from "@/mixin/plugin/config";

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
