import { createPinia } from 'pinia'
import { createPersistedState } from 'pinia-plugin-persistedstate' // 数据持久化
import { useUserStore } from '@/store/user'
import { useGlobalStore } from '@/store/global'
import { useContentStore } from '@/store/content'
import { useWopanStore } from '@/store/wopan'

/* const store = createPinia()
store.use(
  createPersistedState({
    storage: {
      getItem: uni.getStorageSync,
      setItem: uni.setStorageSync,
    },
  }),
) */
