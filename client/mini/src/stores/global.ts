import { defineStore } from "pinia";
import { Platform } from "@/model/config";


interface GlobalState {
  platform: Platform;
}

const state:GlobalState = {
  platform:Platform.Weapp
}

export const useGlobalStore = defineStore({
  id: "global",
  state: () => state,
  getters: {
  },
  actions: {
  },
});
