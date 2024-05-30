import { createPinia } from 'pinia'
import { createPersistedState } from 'pinia-plugin-persistedstate' // 数据持久化
import { useUserStore } from '@/store/user'
import { useGlobalStore } from '@/store/global'
import { useContentStore } from '@/store/content'

const store = createPinia()
store.use(
  createPersistedState({
    storage: {
      getItem: uni.getStorageSync,
      setItem: uni.setStorageSync,
    },
  }),
)

export let userStore
export let globalStore
export let contentStore

export function init() {
  globalStore = useGlobalStore()
  userStore = useUserStore()
  contentStore = useContentStore()
  if (!userStore.auth) {
    userStore.getAuth()
  }
}

export default store

// 模块统一导出
export * from './user'
