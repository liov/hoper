import { createPinia } from 'pinia'
import { createPersistedState } from 'pinia-plugin-persistedstate' // 数据持久化
import { useUserStore } from '@/store/user'
import { useGlobalStore } from '@/store/global'
import { useContentStore } from '@/store/content'
import { useWopanStore } from '@/store/wopan'

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
export let wopanStore

export function init() {
  globalStore = useGlobalStore()
  userStore = useUserStore()
  contentStore = useContentStore()
  wopanStore = useWopanStore()
  if (!userStore.auth) {
    userStore.getAuth()
  }
}

export default store

// 模块统一导出
export * from './user'
