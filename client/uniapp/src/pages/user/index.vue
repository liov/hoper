<template>
  <view class="min-h-screen bg-[#f7f8fa]" :style="{ paddingTop: `${safeAreaInsets?.top || 0}px` }">
    <view class="mx-4 mt-4 rounded-3xl bg-white p-5 shadow-sm">
      <view class="flex items-center gap-3">
        <view class="relative h-15 w-15" @click="onUploadAvatar">
          <image v-if="displayAvatar" class="h-15 w-15 rounded-full bg-[#f3f4f6]" :src="displayAvatar" mode="aspectFill" />
          <view v-else class="h-15 w-15 rounded-full bg-[#f3f4f6] text-[#9ca3af] center text-xs">{{ t('mine.avatarUpload') }}</view>
        </view>
        <view class="text-lg font-semibold text-[#111827]">{{ t('mine.title') }}</view>
      </view>
      <view class="mt-4 text-sm text-[#4b5563]">{{ t('mine.nickname') }}：{{ displayName }}</view>
      <view class="mt-2 text-sm text-[#4b5563]">{{ t('mine.account') }}：{{ displayAccount }}</view>
      <view class="mt-6 flex gap-3">
        <button class="h-10 flex-1 rounded-xl border-none bg-[#018d71] text-white text-sm center" @click="toMoments">{{ t('mine.moments') }}</button>
        <button class="h-10 flex-1 rounded-xl border border-[#fecaca] bg-[#fff1f2] text-[#e11d48] text-sm center active:opacity-85" @click="userStore.logout">{{ t('mine.logout') }}</button>
      </view>
    </view>
  </view>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useUserStore } from '@/store/user'
import { useI18n } from 'vue-i18n'
import { STATIC_DIR as staticDir } from '@/env/config'
import FileService from '@/api/file'
import i18n from '@/locale'

definePage({
  type: 'page',
  style: {
    navigationBarTitleText: i18n.global.t('page.mine'),
  }
})

defineOptions({
  name: 'UserHome',
})

const { safeAreaInsets } = uni.getSystemInfoSync()
const userStore = useUserStore()
const { t } = useI18n()
const localAvatar = ref('')

const displayName = computed(() => userStore.auth?.name || t('mine.defaultName'))
const displayAccount = computed(() => userStore.auth?.mail || userStore.auth?.phone || t('mine.defaultAccount'))
const displayAvatar = computed(() => {
  const avatar = localAvatar.value || userStore.auth?.avatar
  if (!avatar) return ''
  if (avatar.startsWith('http://') || avatar.startsWith('https://')) return avatar
  return `${staticDir}${avatar}`
})

async function onUploadAvatar() {
  try {
    const uploaded = await FileService.chooseAndUploadImage({ count: 1, sizeType: ['compressed'], sourceType: ['album', 'camera'] })
    localAvatar.value = uploaded.fileUrl
    uni.showToast({ title: t('mine.avatarUploaded'), icon: 'success' })
  } catch (e) {
    console.log(e)
    uni.showToast({ title: t('mine.avatarUploadFailed'), icon: 'none' })
  }
}

function toMoments() {
  uni.switchTab({ url: '/pages/moment/moment_list' })
}

onShow(async () => {
  await userStore.getAuth()
  const token = uni.getStorageSync('token')
  if (userStore.auth || token) return
  uni.reLaunch({ url: '/pages/user/login' })
})
</script>

<style scoped>
button::after {
  border: none;
}
</style>
