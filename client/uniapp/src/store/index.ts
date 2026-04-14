import { createPinia } from 'pinia'
import { createPersistedState } from 'pinia-plugin-persistedstate' // 数据持久化

export const store = createPinia()
store.use(
  createPersistedState({
    storage: {
      getItem: uni.getStorageSync,
      setItem: uni.setStorageSync,
    },
  }),
)
