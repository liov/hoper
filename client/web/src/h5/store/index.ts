import { useUserStore } from "@h5/store/user";
import { useGlobalStore } from "@h5/store/global";
import { useContentStore } from "@h5/store/content";

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
