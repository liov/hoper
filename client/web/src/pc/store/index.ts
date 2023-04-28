import { useUserStore } from "@pc/store/user";
import { useGlobalStore } from "@pc/store/global";
import { useContentStore } from "@pc/store/content";

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
