import { useUserStore } from "@/mixin/store/user";
import { useGlobalStore } from "@/mixin/store/global";
import { useContentStore } from "@/mixin/store/content";

export let userStore;
export let globalStore;
export let contentStore;

export function init() {
  globalStore = useGlobalStore();
  userStore = useUserStore();
  contentStore = useContentStore();
  if (!userStore.auth) {
    userStore.getAuth();
  }
}
