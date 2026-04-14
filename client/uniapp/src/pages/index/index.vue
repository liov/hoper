<template>
  <view
    class="bg-white overflow-hidden pt-2 px-4"
    :style="{ marginTop: safeAreaInsets?.top + 'px' }"
  >
    <navigator url="/pages/wopan/list" hover-class="navigator-hover">
      <button type="default">wopan</button>
    </navigator>
  </view>
</template>

<script lang="ts" setup>
import PLATFORM from '@/utils/platform'
import { useToast } from 'wot-design-uni'
import * as wopan from '@hopeio/utils/wopan'

definePage({
  style: {
    navigationBarTitleText: '首页',
  }
})

defineOptions({
  name: 'Home',
})

// 获取屏幕边界到安全区域距离
const { safeAreaInsets } = uni.getSystemInfoSync()
</script>

<style>
.footer {
  padding: 12px;
}
</style>
