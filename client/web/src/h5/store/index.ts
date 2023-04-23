import { useUserStore } from "@/store/user";
import { useGlobalStore } from "@/store/global";
import { useContentStore } from "@/store/content";

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
