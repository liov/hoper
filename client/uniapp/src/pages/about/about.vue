<route lang="json5">
{
  style: {
    navigationBarTitleText: '关于',
  },
}
</route>

<template>
  <view
    class="bg-white overflow-hidden pt-2 px-4"
    :style="{ marginTop: safeAreaInsets?.top + 'px' }"
  >
    <view class="text-center mt-8 text-#fff">
      <wd-button type="success" @click="gotoPage('i18n')">进入多语言页面</wd-button>
    </view>
    <UploadComp />
  </view>
</template>

<script lang="ts" setup>

import UploadComp from './components/upload.vue'

// 获取屏幕边界到安全区域距离
const { safeAreaInsets } = uni.getSystemInfoSync()

const gotoPage = (path) => {
  uni.navigateTo({
    url: `/pages/about/${path}`,
  })
}
</script>

<style lang="scss" scoped>
.test-css {
  // mt-4=>1rem=>16px;
  margin-top: 16px;
}
</style>
