<route lang="json5">
{
  style: {
    navigationStyle: 'custom',
    navigationBarTitleText: 'wopan列表',
  },
}
</route>
<template><view></view></template>

<script lang="ts" setup>
import { useWopanStore } from '@/store/wopan'

defineOptions({
  name: 'WopanList',
})
const wopanStore = useWopanStore()
if (wopanStore.$state.accessToken === '') {
  uni.navigateTo({
    url: '/pages/wopan/login',
  })
}
if (wopanStore.$state.psToken === '') {
  console.log('psToken is empty')
}
</script>

<style scoped></style>
