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
          label="用户名"
          label-width="100px"
          prop="phone"
          clearable
          v-model="model.phone"
          placeholder="请输入用户名"
          :rules="[{ required: true, message: '请填写用户名' }]"
        />
        <wd-input
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

defineOptions({
  name: 'WopanLogin',
})
const wopanStore = useWopanStore()
// 获取屏幕边界到安全区域距离
const { safeAreaInsets } = uni.getSystemInfoSync()
const { success: showSuccess } = useToast()

const model = reactive<{
  phone: string
  smsCode: string
}>({
  phone: '',
  smsCode: '',
})

onLoad(() => {
  console.log('a')
})

const form = ref()

function sendSmsCode() {
  wopan.sendMessageCodeBase(model.phone)
}

function handleSubmit() {
  form.value
    .validate()
    .then(async ({ valid, errors }) => {
      if (valid) {
        await wopanStore.AppLoginByMobile(model)
        showSuccess({
          msg: '校验通过',
        })
      }
    })
    .catch((error) => {
      console.log(error, 'error')
    })
}
</script>

<style>
.footer {
  padding: 12px;
}
</style>
