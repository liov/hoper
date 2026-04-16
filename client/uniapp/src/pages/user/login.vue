<template>
  <view class="min-h-screen flex flex-col bg-gradient-to-br from-[#018d71] via-[#01b892] to-[#e8f8f4]">
    <view :style="{ height: safeAreaInsets?.top + 'px' }" />

    <!-- Logo 区域 -->
    <view class="flex flex-col items-center pt-10 pb-5 flex-shrink-0">
      <image class="w-15 h-15 rounded-[14px] bg-white/20 p-2" src="/static/logo.svg" mode="aspectFit" />
      <text class="text-4xl font-bold text-white mt-5 tracking-widest">Hoper</text>
      <text class="text-sm text-white/80 mt-3">{{ $t('auth.slogan') }}</text>
    </view>

    <!-- #ifdef H5 -->
    <div id="cf-turnstile" class="flex justify-center mt-4" />
    <!-- #endif -->

    <!-- 表单卡片 -->
    <view class="flex-1 bg-white rounded-t-3xl px-6 pb-15 shadow-[0_-4px_20px_rgba(0,0,0,0.08)]">
      <!-- Tab 切换 -->
      <view class="relative flex mx-auto w-45 bg-gray-100 rounded-full p-1.5 mt-6 mb-6 overflow-hidden">
        <view class="flex-1 text-center py-2 text-sm relative z-1 transition-colors duration-200"
          :class="mode === 'login' ? 'text-[#018d71] font-semibold' : 'text-gray-400'" @click="mode = 'login'">{{
            $t('auth.login') }}</view>
        <view class="flex-1 text-center py-2 text-sm relative z-1 transition-colors duration-200"
          :class="mode === 'register' ? 'text-[#018d71] font-semibold' : 'text-gray-400'" @click="mode = 'register'">{{
            $t('auth.register') }}</view>
        <view
          class="absolute top-1.5 left-1.5 bottom-1.5 bg-white rounded-full shadow transition-transform duration-300"
          :style="{ width: 'calc(50% - 6px)', transform: `translateX(${mode === 'login' ? '0%' : '100%'})` }" />
      </view>

      <!-- 登录表单 -->
      <view v-if="mode === 'login'" class="flex flex-col gap-3.5">
        <view class="flex items-center bg-gray-50 rounded-2xl px-3.5 h-13 border-2 border-transparent transition-all"
          :class="focused === 'account' ? 'border-[#018d71] bg-[#f0fdf8]' : ''">
          <text class="text-xl mr-2.5 flex-shrink-0">📱</text>
          <input v-model="loginForm.account" class="flex-1 text-sm text-gray-800 bg-transparent h-full" type="text"
            :placeholder="$t('auth.account')" placeholder-class="text-gray-300 text-sm" @focus="focused = 'account'"
            @blur="focused = ''" />
        </view>
        <view class="flex items-center bg-gray-50 rounded-2xl px-3.5 h-13 border-2 border-transparent transition-all"
          :class="focused === 'loginPwd' ? 'border-[#018d71] bg-[#f0fdf8]' : ''">
          <text class="text-xl mr-2.5 flex-shrink-0">🔒</text>
          <input v-model="loginForm.password" class="flex-1 text-sm text-gray-800 bg-transparent h-full"
            :class="{ 'pwd-masked': !showLoginPwd && loginForm.password.length > 0 }" type="text"
            :placeholder="$t('auth.password')"
            placeholder-class="text-gray-300 text-sm" @focus="focused = 'loginPwd'" @blur="focused = ''" />
          <text class="ml-2 flex-shrink-0 text-gray-400 text-base select-none" @click.stop="showLoginPwd = !showLoginPwd">{{ showLoginPwd ? '🙈' : '👁' }}</text>
        </view>

        <view class="flex justify-end -mt-2">
          <text class="text-xs text-[#018d71]" @click="onForgotPwd">{{ $t('auth.forgotPwd') }}</text>
        </view>

        <button
          class="w-full h-12 bg-gradient-to-r from-[#018d71] to-[#01b892] text-white text-base font-semibold rounded-full border-none tracking-widest mt-2 shadow-[0_4px_12px_rgba(1,141,113,0.35)] center"
          :disabled="submitting" @click="onLogin">
          <text v-if="!submitting">{{ $t('auth.loginBtn') }}</text>
          <view v-else class="flex gap-1.5 items-center">
            <view class="dot" />
            <view class="dot" />
            <view class="dot" />
          </view>
        </button>

        <view class="flex items-center gap-2.5 my-1">
          <view class="flex-1 h-px bg-gray-100" />
          <text class="text-xs text-gray-300 whitespace-nowrap">{{ $t('auth.otherLogin') }}</text>
          <view class="flex-1 h-px bg-gray-100" />
        </view>
        <view class="flex justify-center gap-15">
          <view class="flex flex-col items-center gap-1.5" @click="onThirdLogin('wechat')">
            <view class="w-12 h-12 rounded-full bg-gray-50 center border border-gray-100 text-xl">💬</view>
            <text class="text-xs text-gray-400">{{ $t('auth.wechat') }}</text>
          </view>
          <view class="flex flex-col items-center gap-1.5" @click="onThirdLogin('apple')">
            <view class="w-12 h-12 rounded-full bg-gray-50 center border border-gray-100 text-xl">🍎</view>
            <text class="text-xs text-gray-400">{{ $t('auth.apple') }}</text>
          </view>
        </view>
      </view>

      <!-- 注册表单 -->
      <view v-else class="flex flex-col gap-3.5">
        <view class="flex items-center bg-gray-50 rounded-2xl px-3.5 h-13 border-2 border-transparent transition-all"
          :class="focused === 'phone' ? 'border-[#018d71] bg-[#f0fdf8]' : ''">
          <text class="text-xl mr-2.5 flex-shrink-0">📱</text>
          <input v-model="registerForm.phone" class="flex-1 text-sm text-gray-800 bg-transparent h-full" type="number"
            maxlength="11" :placeholder="$t('auth.phone')" placeholder-class="text-gray-300 text-sm"
            @focus="focused = 'phone'" @blur="focused = ''" />
        </view>
        <view class="flex items-center bg-gray-50 rounded-2xl px-3.5 h-13 border-2 border-transparent transition-all"
          :class="focused === 'sms' ? 'border-[#018d71] bg-[#f0fdf8]' : ''">
          <text class="text-xl mr-2.5 flex-shrink-0">✉️</text>
          <input v-model="registerForm.smsCode" class="flex-1 text-sm text-gray-800 bg-transparent h-full" type="number"
            maxlength="6" :placeholder="$t('auth.smsCode')" placeholder-class="text-gray-300 text-sm"
            @focus="focused = 'sms'" @blur="focused = ''" />
          <button class="text-xs border-none rounded-lg px-2.5 py-1.5 ml-2 flex-shrink-0"
            :class="smsCountdown > 0 ? 'text-gray-400 bg-gray-100' : 'text-[#018d71] bg-[#e6f7f3]'"
            :disabled="smsCountdown > 0" @click="onSendSms">{{ smsCountdown > 0 ? `${smsCountdown}s` : $t('auth.getSms')
            }}</button>
        </view>
        <view class="flex items-center bg-gray-50 rounded-2xl px-3.5 h-13 border-2 border-transparent transition-all"
          :class="focused === 'regPwd' ? 'border-[#018d71] bg-[#f0fdf8]' : ''">
          <text class="text-xl mr-2.5 flex-shrink-0">🔒</text>
          <input v-model="registerForm.password" class="flex-1 text-sm text-gray-800 bg-transparent h-full"
            :class="{ 'pwd-masked': !showRegPwd && registerForm.password.length > 0 }" type="text"
            :placeholder="$t('auth.setPwd')"
            placeholder-class="text-gray-300 text-sm" @focus="focused = 'regPwd'" @blur="focused = ''" />
          <text class="ml-2 flex-shrink-0 text-gray-400 text-base select-none" @click.stop="showRegPwd = !showRegPwd">{{ showRegPwd ? '🙈' : '👁' }}</text>
        </view>
        <view class="flex items-center bg-gray-50 rounded-2xl px-3.5 h-13 border-2 border-transparent transition-all"
          :class="focused === 'confirm' ? 'border-[#018d71] bg-[#f0fdf8]' : ''">
          <text class="text-xl mr-2.5 flex-shrink-0">🔒</text>
          <input v-model="registerForm.confirmPassword" class="flex-1 text-sm text-gray-800 bg-transparent h-full"
            :class="{ 'pwd-masked': !showRegPwd && registerForm.confirmPassword.length > 0 }" type="text"
            :placeholder="$t('auth.confirmPwd')"
            placeholder-class="text-gray-300 text-sm" @focus="focused = 'confirm'" @blur="focused = ''" />
        </view>

        <view class="flex items-center flex-wrap gap-1 -mt-1">
          <view class="w-4.5 h-4.5 rounded flex-shrink-0 border-2 center transition-all mr-1"
            :class="agreed ? 'bg-[#018d71] border-[#018d71]' : 'border-gray-300'" @click="agreed = !agreed"><text
              v-if="agreed" class="text-xs text-white font-bold">✓</text></view>
          <text class="text-xs text-gray-400">{{ $t('auth.agreePrefix') }}</text>
          <text class="text-xs text-[#018d71]" @click="onViewAgreement('user')">{{ $t('auth.userAgreement') }}</text>
          <text class="text-xs text-gray-400">{{ $t('auth.and') }}</text>
          <text class="text-xs text-[#018d71]" @click="onViewAgreement('privacy')">{{ $t('auth.privacyPolicy') }}</text>
        </view>

        <button
          class="w-full h-12 bg-gradient-to-r from-[#018d71] to-[#01b892] text-white text-base font-semibold rounded-full border-none tracking-widest mt-2 shadow-[0_4px_12px_rgba(1,141,113,0.35)] center"
          :disabled="submitting" @click="onRegister">
          <text v-if="!submitting">{{ $t('auth.registerBtn') }}</text>
          <view v-else class="flex gap-1.5 items-center">
            <view class="dot" />
            <view class="dot" />
            <view class="dot" />
          </view>
        </button>
      </view>
    </view>
  </view>
</template>

<script lang="ts" setup>
import { useUserStore } from '@/store/user'
import { useI18n } from 'vue-i18n'
// #ifdef H5
// #endif

definePage({
  type: 'home',
  style: { navigationStyle: 'custom', navigationBarTitleText: '登录' },
})

const userStore = useUserStore()
const { t } = useI18n()
const { safeAreaInsets } = uni.getSystemInfoSync()

const mode = ref<'login' | 'register'>('login')
const focused = ref('')
const submitting = ref(false)
const showLoginPwd = ref(false)
const showRegPwd = ref(false)
const agreed = ref(false)
const smsCountdown = ref(0)

const loginForm = reactive({ account: '', password: '' })
const registerForm = reactive({ phone: '', smsCode: '', password: '', confirmPassword: '' })

let smsTimer: ReturnType<typeof setInterval> | null = null

function startSmsCountdown() {
  smsCountdown.value = 60
  smsTimer = setInterval(() => {
    if (--smsCountdown.value <= 0 && smsTimer) { clearInterval(smsTimer); smsTimer = null }
  }, 1000)
}

async function onSendSms() {
  if (registerForm.phone.length !== 11) return uni.showToast({ title: t('auth.err.phone'), icon: 'none' })
  startSmsCountdown()
  uni.showToast({ title: t('auth.err.smsSent'), icon: 'success' })
}

async function onLogin() {
  if (!loginForm.account) return uni.showToast({ title: t('auth.err.account'), icon: 'none' })
  if (!loginForm.password) return uni.showToast({ title: t('auth.err.password'), icon: 'none' })
  submitting.value = true
  try { await userStore.login({ account: loginForm.account, password: loginForm.password }) }
  finally { submitting.value = false }
}

async function onRegister() {
  if (registerForm.phone.length !== 11) return uni.showToast({ title: t('auth.err.phone'), icon: 'none' })
  if (!registerForm.smsCode) return uni.showToast({ title: t('auth.err.smsCode'), icon: 'none' })
  if (registerForm.password.length < 8) return uni.showToast({ title: t('auth.err.pwdLength'), icon: 'none' })
  if (registerForm.password !== registerForm.confirmPassword) return uni.showToast({ title: t('auth.err.pwdNotMatch'), icon: 'none' })
  if (!agreed.value) return uni.showToast({ title: t('auth.err.agreement'), icon: 'none' })
  submitting.value = true
  try { await userStore.signup({ phone: registerForm.phone, smsCode: registerForm.smsCode, password: registerForm.password }) }
  finally { submitting.value = false }
}

function onForgotPwd() { uni.showToast({ title: t('auth.err.forgotPwd'), icon: 'none' }) }
function onThirdLogin(type: string) { uni.showToast({ title: t('auth.err.thirdLogin', { type }), icon: 'none' }) }
function onViewAgreement(type: string) { uni.showToast({ title: t(type === 'user' ? 'auth.userAgreement' : 'auth.privacyPolicy'), icon: 'none' }) }

onUnmounted(() => { if (smsTimer) clearInterval(smsTimer) })

// #ifdef H5
onMounted(() => {
  const render = () => {
    (window as any).turnstile.render('#cf-turnstile', {
      sitekey: '0x4AAAAAAAgC9s4WZMlljGRg',
      callback: (token: string) => console.log(`Challenge Success ${token}`),
    })
  }
    ; (window as any).turnstile ? render() : (window as any).onloadTurnstileCallback = render
})
// #endif
</script>

<style scoped>
/* loading 动画无法用 utility 表达 */
.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #fff;
  animation: pulse 1.2s infinite ease-in-out;
}

.dot:nth-child(2) {
  animation-delay: 0.2s;
}

.dot:nth-child(3) {
  animation-delay: 0.4s;
}

.pwd-masked {
  -webkit-text-security: disc;
}

@keyframes pulse {

  0%,
  80%,
  100% {
    transform: scale(0.6);
    opacity: 0.5;
  }

  40% {
    transform: scale(1);
    opacity: 1;
  }
}
</style>
