import { defineStore } from "pinia";
import { Platform } from "@/model/const";

interface GlobalState {
  counter: number;
  platform: Platform;
}

const state: GlobalState = {
  counter: 0,
  platform: Platform.H5,
};

export const useGlobalStore = defineStore({
  id: "global",
  state: () => state,
  getters: {
    doubleCount: (state) => state.counter * 2,
  },
  actions: {
    increment() {
      this.counter++;
    },
    setPlatform(platform: Platform) {
      this.platform = platform;
    },
  },
});
