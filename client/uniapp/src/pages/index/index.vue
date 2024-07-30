<!-- 使用 type="home" 属性设置首页，其他页面不需要设置，默认为page；推荐使用json5，更强大，且允许注释 -->
<route lang="json5" type="home">
{
  style: {
    navigationStyle: 'custom',
    navigationBarTitleText: '首页',
  },
}
</route>
<template>
  <view
    class="bg-white overflow-hidden pt-2 px-4"
    :style="{ marginTop: safeAreaInsets?.top + 'px' }"
  >
    <view class="mt-12">
      <image src="/static/logo.svg" alt="" class="w-28 h-28 block mx-auto" />
    </view>
    <view class="text-center text-4xl main-title-color mt-4">hoper</view>
    <view class="text-center text-2xl mt-2 mb-8">hoper with record</view>

    <view class="text-justify max-w-100 m-auto text-4 indent mb-2">{{ description }}</view>
    <navigator url="/pages/moment/moment_list?title=navigate" hover-class="navigator-hover">
      <image class="logo" src="/static/logo.svg" />
    </navigator>
    <view class="text-area">
      <text class="title">{{ title }}</text>
    </view>
    <view class="text-center mt-8">
      当前平台是：
      <text class="text-green-500">{{ PLATFORM.platform }}</text>
    </view>
    <view class="text-center mt-4">
      模板分支是：
      <text class="text-green-500">i18n</text>
    </view>
  </view>
  <tabbar />
</template>

<script lang="ts" setup>
import PLATFORM from '@/utils/platform'

defineOptions({
  name: 'Home',
})

// 获取屏幕边界到安全区域距离
const { safeAreaInsets } = uni.getSystemInfoSync()

const description = ref('首页')

onLoad(() => {
  console.log('a')
})
</script>

<style>
.main-title-color {
  color: #d14328;
}
</style>
