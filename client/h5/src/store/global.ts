import { defineStore } from "pinia";

interface GlobalState {
  counter: number;
}

export const useGlobalStore = defineStore({
  id: "global",
  state: () =>
    ({
      counter: 0,
    } as GlobalState),
  getters: {
    doubleCount: (state) => state.counter * 2,
  },
  actions: {
    increment() {
      this.counter++;
    },
  },
});
