<route lang="json5">
{
  style: {
    navigationStyle: 'custom',
    navigationBarTitleText: 'wopan登录',
  },
}
</route>
<template>
  <view
    class="bg-white overflow-hidden pt-2 px-4"
    :style="{ marginTop: safeAreaInsets?.top + 'px' }"
  >
    <wd-form ref="form" :model="model">
      <wd-cell-group border>
        <wd-input
          v-if="!setPwd"
          label="用户名"
          label-width="100px"
          prop="phone"
          clearable
          v-model="model.phone"
          placeholder="请输入用户名"
          :rules="[{ required: true, message: '请填写用户名' }]"
        />
        <wd-input
          v-if="!setPwd"
          label="验证码"
          label-width="100px"
          prop="smsCode"
          clearable
          use-suffix-slot
          v-model="model.smsCode"
          placeholder="请输入验证码"
          :rules="[{ required: true, message: '请填写验证码' }]"
        >
          <template #suffix>
            <view><wd-button @click="sendSmsCode">获取验证码</wd-button></view>
          </template>
        </wd-input>
        <wd-input
          v-if="setPwd"
          label="密码"
          label-width="100px"
          prop="passwd"
          clearable
          v-model="model.passwd"
          placeholder="请输入密码"
          :rules="[{ required: true, message: '请填写密码' }]"
        />
      </wd-cell-group>
      <view class="footer">
        <wd-button type="primary" size="large" @click="handleSubmit" block>提交</wd-button>
      </view>
    </wd-form>
  </view>
</template>

<script lang="ts" setup>
import PLATFORM from '@/utils/platform'
import { useToast } from 'wot-design-uni'
import * as wopan from 'diamond/wopan'

import { useWopanStore } from '@/store/wopan'
import { onLoad } from '@dcloudio/uni-app'

defineOptions({
  name: 'WopanLogin',
})
const setPwd = ref(false)
onLoad((options) => {
  setPwd.value = options.psToken === '1'
})

const wopanStore = useWopanStore()
// 获取屏幕边界到安全区域距离
const { safeAreaInsets } = uni.getSystemInfoSync()
const { success: showSuccess } = useToast()

const model = reactive<{
  phone: string
  smsCode: string
  passwd: string
}>({
  phone: '',
  smsCode: '',
  passwd: '',
})

onLoad(() => {
  console.log('a')
})

const form = ref()

function sendSmsCode() {
  wopan.sendMessageCodeBase(model.phone)
}

async function handleSubmit() {
  if (setPwd.value) {
    await wopanStore.PrivateSpaceLogin(model)
    uni.navigateBack()
  } else {
    await wopanStore.AppLoginByMobile(model)
    setPwd.value = true
  }
}
</script>

<style>
.footer {
  padding: 12px;
}
</style>
