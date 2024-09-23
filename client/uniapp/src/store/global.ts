import { defineStore } from 'pinia'
import { Platform } from '@/env/platform'

export interface GlobalState {
  counter: number
  platform: Platform
}

const state: GlobalState = {
  counter: 0,
  platform: Platform.H5,
}

export const useGlobalStore = defineStore('global', {
  state: () => state,
  getters: {
    doubleCount: (state) => state.counter * 2,
  },
  actions: {
    increment() {
      this.counter++
    },
    setPlatform(platform: Platform) {
      this.platform = platform
    },
  },
})
