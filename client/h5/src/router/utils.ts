import { contentRoute } from "@/router/enum";
import { contentMutations } from "@/store/enum";
import store from "@/store/index";
import router from "@/router/index";
import emitter from "@/plugin/emitter";

export const jump = (path: string, type: number, content: any) => {
  const route = `/${contentRoute[type]}/${content.id}`;
  if (path !== route) {
    store.commit(contentMutations[type], content);
    router.push(route);
  }
  emitter.emit("onComment");
};
