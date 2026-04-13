<template>
  <div class="login-page">
    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
    </div>

    <div class="login-container">
      <!-- 左侧装饰面板 -->
      <div class="brand-panel">
        <div class="brand-content">
          <div class="brand-logo">
            <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="24" cy="24" r="22" stroke="white" stroke-width="2.5" fill="rgba(255,255,255,0.15)"/>
              <path d="M14 24C14 18.477 18.477 14 24 14s10 4.477 10 10-4.477 10-10 10" stroke="white" stroke-width="2.5" stroke-linecap="round"/>
              <circle cx="24" cy="24" r="4" fill="white"/>
            </svg>
          </div>
          <h1 class="brand-name">Hoper</h1>
          <p class="brand-slogan">{{ t('login.pureBrandSlogan') }}</p>
          <div class="brand-features">
            <div class="feature-item">
              <span class="feature-icon">✦</span>
              <span>{{ t('login.pureFeature1') }}</span>
            </div>
            <div class="feature-item">
              <span class="feature-icon">✦</span>
              <span>{{ t('login.pureFeature2') }}</span>
            </div>
            <div class="feature-item">
              <span class="feature-icon">✦</span>
              <span>{{ t('login.pureFeature3') }}</span>
            </div>
          </div>
        </div>
        <div class="brand-wave">
          <svg viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
            <path fill="rgba(255,255,255,0.08)"
              d="M44.7,-76.4C58.8,-69.2,71.8,-59.1,79.6,-45.8C87.4,-32.6,90,-16.3,88.5,-1C87,14.3,81.5,28.5,73.5,41.3C65.5,54,55,65.1,42.2,72.4C29.5,79.6,14.7,83,-0.9,84.4C-16.5,85.8,-33,85.2,-46.6,78.3C-60.2,71.4,-70.9,58.2,-77.6,43.3C-84.2,28.4,-86.8,11.8,-84.9,-3.8C-83,-19.4,-76.6,-34,-67.2,-46.2C-57.8,-58.4,-45.4,-68.1,-32,-75.1C-18.6,-82.1,-4.2,-86.4,9.7,-84.6C23.6,-82.8,30.7,-83.5,44.7,-76.4Z"
              transform="translate(100 100)" />
          </svg>
        </div>
      </div>

      <!-- 右侧表单面板 -->
      <div class="form-panel">
        <div class="form-header">
          <h2 class="form-title">{{ isRegister ? t('login.pureCreateAccount') : t('login.pureWelcomeBack') }}</h2>
          <p class="form-subtitle">{{ isRegister ? t('login.pureRegisterSubtitle') : t('login.pureLoginSubtitle') }}</p>
        </div>

        <!-- 切换标签 -->
        <div class="tab-switcher">
          <button class="tab-btn" :class="{ active: !isRegister }" @click="switchTab(false)">{{ t('login.pureLogin') }}</button>
          <button class="tab-btn" :class="{ active: isRegister }" @click="switchTab(true)">{{ t('login.pureRegister') }}</button>
          <div class="tab-indicator" :style="{ transform: isRegister ? 'translateX(100%)' : 'translateX(0)' }"></div>
        </div>

        <!-- 登录表单 -->
        <Transition name="slide-fade" mode="out-in">
          <form v-if="!isRegister" key="login" class="auth-form" @submit.prevent="handleSubmit">
            <div class="field-group">
              <label class="field-label">{{ t('login.pureAccount') }}</label>
              <div class="input-wrapper" :class="{ focused: focusedField === 'account', error: errors.account }">
                <span class="input-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                  </svg>
                </span>
                <input v-model="account" type="text" :placeholder="t('login.pureAccountPlaceholder')"
                  @focus="focusedField = 'account'" @blur="focusedField = ''" />
              </div>
              <span v-if="errors.account" class="error-msg">{{ errors.account }}</span>
            </div>

            <div class="field-group">
              <label class="field-label">
                {{ t('login.purePassword') }}
                <a class="forgot-link" href="javascript:void(0)">{{ t('login.pureForget') }}</a>
              </label>
              <div class="input-wrapper" :class="{ focused: focusedField === 'password', error: errors.password }">
                <span class="input-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                    <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                  </svg>
                </span>
                <input v-model="password" :type="showPassword ? 'text' : 'password'" :placeholder="t('login.purePasswordPlaceholder')"
                  @focus="focusedField = 'password'" @blur="focusedField = ''" />
                <button type="button" class="eye-btn" @click="showPassword = !showPassword">
                  <svg v-if="!showPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  </svg>
                  <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                    <line x1="1" y1="1" x2="23" y2="23"/>
                  </svg>
                </button>
              </div>
              <span v-if="errors.password" class="error-msg">{{ errors.password }}</span>
            </div>

            <div class="captcha-wrapper">
              <Luosimao ref="luosimao" />
            </div>

            <button type="submit" class="submit-btn" :disabled="loading">
              <span v-if="loading" class="loading-spinner"></span>
              <span>{{ loading ? t('login.pureLoginLoading') : t('login.pureLoginBtn') }}</span>
            </button>

            <div class="divider"><span>{{ t('login.pureOr') }}</span></div>

            <div class="social-login">
              <button type="button" class="social-btn wechat" :title="t('login.pureWeChatLogin')">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M8.07 15.29c-.38 0-.73-.14-1-.38l-3.15-2.63.72-.86 2.93 2.44 6.28-5.24.72.86-6.28 5.24c-.27.24-.62.38-1 .37zm13.35-3.38c0 2.67-1.9 4.93-4.56 5.87L16 20l-2.52-1.46a7.6 7.6 0 0 1-1.73.2c-3.73 0-6.75-2.69-6.75-6.01s3.02-6.01 6.75-6.01S21.42 9.4 21.42 12z"/>
                </svg>
                {{ t('login.pureWeChatLogin') }}
              </button>
            </div>
          </form>

          <!-- 注册表单 -->
          <form v-else key="register" class="auth-form" @submit.prevent="handleSubmit">
            <div class="field-group">
              <label class="field-label">{{ t('login.pureName') }}</label>
              <div class="input-wrapper" :class="{ focused: focusedField === 'username', error: errors.username }">
                <span class="input-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                  </svg>
                </span>
                <input v-model="username" type="text" :placeholder="t('login.pureUsernamePlaceholder')"
                  @focus="focusedField = 'username'" @blur="focusedField = ''" />
              </div>
              <span v-if="errors.username" class="error-msg">{{ errors.username }}</span>
            </div>

            <div class="field-group">
              <label class="field-label">{{ t('login.pureGender') }}</label>
              <div class="gender-group">
                <label class="gender-option" :class="{ selected: gender === '1' }">
                  <input type="radio" v-model="gender" value="1" hidden />
                  <span class="gender-icon">♂</span>
                  <span>{{ t('login.pureMale') }}</span>
                </label>
                <label class="gender-option" :class="{ selected: gender === '2' }">
                  <input type="radio" v-model="gender" value="2" hidden />
                  <span class="gender-icon">♀</span>
                  <span>{{ t('login.pureFemale') }}</span>
                </label>
              </div>
            </div>

            <div class="fields-row">
              <div class="field-group">
                <label class="field-label">{{ t('login.purePhone') }}</label>
                <div class="input-wrapper" :class="{ focused: focusedField === 'phone', error: errors.phone }">
                  <span class="input-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07A19.5 19.5 0 0 1 4.69 12a19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 3.6 1.23h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L7.91 8.57a16 16 0 0 0 6 6l.18-.18a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 21 16.92z"/>
                    </svg>
                  </span>
                  <input v-model="phone" type="tel" :placeholder="t('login.purePhonePlaceholder')"
                    @focus="focusedField = 'phone'" @blur="focusedField = ''" />
                </div>
                <span v-if="errors.phone" class="error-msg">{{ errors.phone }}</span>
              </div>

              <div class="field-group">
                <label class="field-label">{{ t('login.pureEmail') }}</label>
                <div class="input-wrapper" :class="{ focused: focusedField === 'mail', error: errors.mail }">
                  <span class="input-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                      <polyline points="22,6 12,13 2,6"/>
                    </svg>
                  </span>
                  <input v-model="mail" type="email" :placeholder="t('login.pureEmailPlaceholder')"
                    @focus="focusedField = 'mail'" @blur="focusedField = ''" />
                </div>
                <span v-if="errors.mail" class="error-msg">{{ errors.mail }}</span>
              </div>
            </div>

            <div class="field-group">
              <label class="field-label">{{ t('login.purePassword') }}</label>
              <div class="input-wrapper" :class="{ focused: focusedField === 'reg-password', error: errors.password }">
                <span class="input-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                    <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                  </svg>
                </span>
                <input v-model="password" :type="showPassword ? 'text' : 'password'" :placeholder="t('login.purePasswordMinPlaceholder')"
                  @focus="focusedField = 'reg-password'" @blur="focusedField = ''" />
                <button type="button" class="eye-btn" @click="showPassword = !showPassword">
                  <svg v-if="!showPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  </svg>
                  <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                    <line x1="1" y1="1" x2="23" y2="23"/>
                  </svg>
                </button>
              </div>
              <div v-if="password" class="password-strength">
                <div class="strength-bars">
                  <div class="bar" :class="{ active: pwStrength >= 1, weak: pwStrength === 1, medium: pwStrength === 2, strong: pwStrength >= 3 }"></div>
                  <div class="bar" :class="{ active: pwStrength >= 2, medium: pwStrength === 2, strong: pwStrength >= 3 }"></div>
                  <div class="bar" :class="{ active: pwStrength >= 3, strong: pwStrength >= 3 }"></div>
                </div>
                <span class="strength-label">{{ pwStrengthLabel ? t(`login.${pwStrengthLabel}`) : '' }}</span>
              </div>
              <span v-if="errors.password" class="error-msg">{{ errors.password }}</span>
            </div>

            <div class="field-group">
              <label class="field-label">{{ t('login.pureSure') }}</label>
              <div class="input-wrapper" :class="{ focused: focusedField === 'confirm', error: errors.passwordConfirm }">
                <span class="input-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                  </svg>
                </span>
                <input v-model="passwordConfirm" :type="showPassword ? 'text' : 'password'" :placeholder="t('login.purePasswordConfirmPlaceholder')"
                  @focus="focusedField = 'confirm'" @blur="focusedField = ''" />
              </div>
              <span v-if="errors.passwordConfirm" class="error-msg">{{ errors.passwordConfirm }}</span>
            </div>

            <div class="captcha-wrapper">
              <Luosimao ref="luosimao" />
            </div>

            <button type="submit" class="submit-btn" :disabled="loading">
              <span v-if="loading" class="loading-spinner"></span>
              <span>{{ loading ? t('login.pureRegisterLoading') : t('login.pureRegisterBtn') }}</span>
            </button>
          </form>
        </Transition>

        <p class="switch-hint">
          {{ isRegister ? t('login.pureHasAccount') : t('login.pureNoAccount') }}
          <a href="javascript:void(0)" class="switch-link" @click="switchTab(!isRegister)">
            {{ isRegister ? t('login.pureGoLogin') : t('login.pureGoRegister') }}
          </a>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import Luosimao from "@/components/Luosimao.vue";
import { Validator } from "@hopeio/utils/validator";
import { useRoute, useRouter } from "vue-router";
import { useUserStore } from "@/store/user";
import { Dialog, Toast } from "vant";
import { useAppStore } from "@/store/modules/app";
import { Platform } from "@/types/enum";

const { t } = useI18n();

const appState = useAppStore();
const router = useRouter();
const route = useRoute();
const store = useUserStore();
const luosimao = ref<any>(null);

const isRegister = ref(false);
const loading = ref(false);
const focusedField = ref("");
const showPassword = ref(false);

const account = ref("");
const username = ref("");
const password = ref("");
const passwordConfirm = ref("");
const gender = ref("1");
const mail = ref("");
const phone = ref("");

const errors = ref<Record<string, string>>({});

const pwStrength = computed(() => {
  const p = password.value;
  if (!p) return 0;
  let score = 0;
  if (p.length >= 8) score++;
  if (/[A-Z]/.test(p) || /[0-9]/.test(p)) score++;
  if (/[^A-Za-z0-9]/.test(p)) score++;
  return score;
});

const pwStrengthLabel = computed(() => {
  const keys = ["", "purePasswordWeak", "purePasswordMedium", "purePasswordStrong"];
  return keys[pwStrength.value] ?? "";
});

onMounted(() => {
  if (appState.platform === Platform.Weapp) {
    Dialog.confirm({ title: t('login.pureWeChatLogin'), message: t('login.pureWeChatLogin') })
      .then(() => {
        window.wx.miniProgram.navigateTo({
          url: `/pages/user/login?h5Url=${encodeURIComponent(route.fullPath)}`,
        });
      })
      .catch(() => {});
  }
});

if (route.query.back) {
  if (!store.auth) await store.getAuth();
  if (store.auth) await router.replace(`${route.query.back}`);
}

function switchTab(toRegister: boolean) {
  isRegister.value = toRegister;
  errors.value = {};
  password.value = "";
  showPassword.value = false;
}

function validate(): boolean {
  const errs: Record<string, string> = {};
  if (!isRegister.value) {
    if (!account.value) errs.account = t('login.pureAccountReg');
    if (!password.value) errs.password = t('login.purePassWordReg');
  } else {
    if (!username.value) errs.username = t('login.pureUsernameReg');
    if (!phone.value || !Validator.PhoneReg.test(phone.value)) errs.phone = t('login.purePhoneCorrectReg');
    if (!mail.value || !Validator.EmailReg.test(mail.value)) errs.mail = t('login.pureEmailReg');
    if (!password.value || password.value.length < 8) errs.password = t('login.purePasswordMinReg');
    if (passwordConfirm.value !== password.value) errs.passwordConfirm = t('login.purePassWordDifferentReg');
  }
  errors.value = errs;
  return Object.keys(errs).length === 0;
}

async function handleSubmit() {
  if (!validate()) return;
  loading.value = true;
  try {
    const vCode = luosimao.value?.getValue?.() ?? "";
    if (!isRegister.value) {
      await store.login({ input: account.value, password: password.value, vCode });
    } else {
      await store.signup({ name: username.value, gender: parseInt(gender.value), phone: phone.value, mail: mail.value, password: password.value, vCode });
      Toast.success(t('login.pureCheckEmail'));
    }
  } finally {
    loading.value = false;
    const LUOCAPTCHA = window.LUOCAPTCHA;
    LUOCAPTCHA && LUOCAPTCHA.reset();
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
  padding: 20px;
  overflow: hidden;
}

/* 背景动态圆圈 */
.bg-decoration { position: fixed; inset: 0; pointer-events: none; overflow: hidden; }
.circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.06);
  animation: float 8s ease-in-out infinite;
}
.circle-1 { width: 400px; height: 400px; top: -100px; right: -100px; animation-delay: 0s; }
.circle-2 { width: 250px; height: 250px; bottom: -50px; left: -50px; animation-delay: -3s; }
.circle-3 { width: 180px; height: 180px; top: 50%; left: 40%; animation-delay: -5s; }

@keyframes float {
  0%, 100% { transform: translateY(0) scale(1); }
  50% { transform: translateY(-30px) scale(1.05); }
}

/* 主容器 */
.login-container {
  display: flex;
  width: 100%;
  max-width: 900px;
  min-height: 580px;
  border-radius: 24px;
  overflow: hidden;
  box-shadow: 0 40px 80px rgba(0, 0, 0, 0.3), 0 0 0 1px rgba(255,255,255,0.1);
  backdrop-filter: blur(20px);
}

/* 左侧品牌面板 */
.brand-panel {
  flex: 0 0 380px;
  background: linear-gradient(160deg, rgba(102,126,234,0.9) 0%, rgba(118,75,162,0.95) 100%);
  padding: 48px 40px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  position: relative;
  overflow: hidden;
}

.brand-content { position: relative; z-index: 1; }

.brand-logo svg {
  width: 56px;
  height: 56px;
  margin-bottom: 20px;
  filter: drop-shadow(0 4px 12px rgba(0,0,0,0.2));
}

.brand-name {
  font-size: 36px;
  font-weight: 700;
  color: white;
  letter-spacing: 2px;
  margin-bottom: 8px;
}

.brand-slogan {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.75);
  margin-bottom: 48px;
  letter-spacing: 1px;
}

.brand-features { display: flex; flex-direction: column; gap: 18px; }

.feature-item {
  display: flex;
  align-items: center;
  gap: 12px;
  color: rgba(255, 255, 255, 0.88);
  font-size: 14px;
}

.feature-icon {
  width: 28px;
  height: 28px;
  background: rgba(255,255,255,0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  color: white;
  flex-shrink: 0;
}

.brand-wave {
  position: absolute;
  bottom: -40px;
  right: -40px;
  width: 240px;
  height: 240px;
  opacity: 0.5;
}

/* 右侧表单面板 */
.form-panel {
  flex: 1;
  background: rgba(255, 255, 255, 0.97);
  padding: 44px 48px;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.form-header { margin-bottom: 28px; }

.form-title {
  font-size: 26px;
  font-weight: 700;
  color: #1a1a2e;
  margin-bottom: 6px;
}

.form-subtitle {
  font-size: 14px;
  color: #6b7280;
}

/* Tab 切换 */
.tab-switcher {
  display: flex;
  position: relative;
  background: #f3f4f6;
  border-radius: 10px;
  padding: 4px;
  margin-bottom: 28px;
  width: 200px;
}

.tab-btn {
  flex: 1;
  padding: 8px 0;
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
  background: transparent;
  border: none;
  cursor: pointer;
  border-radius: 8px;
  position: relative;
  z-index: 1;
  transition: color 0.25s;
}

.tab-btn.active { color: #4f46e5; }

.tab-indicator {
  position: absolute;
  top: 4px;
  left: 4px;
  width: calc(50% - 4px);
  height: calc(100% - 8px);
  background: white;
  border-radius: 7px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.12);
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

/* 表单字段 */
.auth-form { display: flex; flex-direction: column; gap: 16px; flex: 1; }

.field-group { display: flex; flex-direction: column; gap: 6px; }

.field-label {
  font-size: 13px;
  font-weight: 600;
  color: #374151;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.forgot-link {
  font-size: 12px;
  color: #6366f1;
  text-decoration: none;
  font-weight: 400;
}
.forgot-link:hover { text-decoration: underline; }

.input-wrapper {
  display: flex;
  align-items: center;
  gap: 10px;
  border: 1.5px solid #e5e7eb;
  border-radius: 10px;
  padding: 0 14px;
  background: #fafafa;
  transition: border-color 0.2s, box-shadow 0.2s, background 0.2s;
}

.input-wrapper.focused {
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.12);
  background: white;
}

.input-wrapper.error { border-color: #ef4444; box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1); }

.input-icon {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}
.input-icon svg { width: 16px; height: 16px; color: #9ca3af; }

.input-wrapper input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  padding: 13px 0;
  font-size: 14px;
  color: #111827;
}

.input-wrapper input::placeholder { color: #d1d5db; }

.eye-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  color: #9ca3af;
  display: flex;
  align-items: center;
  transition: color 0.2s;
}
.eye-btn:hover { color: #6366f1; }
.eye-btn svg { width: 16px; height: 16px; }

.error-msg { font-size: 12px; color: #ef4444; }

/* 性别选择 */
.gender-group { display: flex; gap: 12px; }

.gender-option {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px;
  border: 1.5px solid #e5e7eb;
  border-radius: 10px;
  cursor: pointer;
  font-size: 14px;
  color: #6b7280;
  transition: all 0.2s;
  background: #fafafa;
}

.gender-option.selected {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.06);
  color: #4f46e5;
  font-weight: 600;
}

.gender-icon { font-size: 16px; }

/* 两列布局 */
.fields-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }

/* 密码强度 */
.password-strength {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
}

.strength-bars { display: flex; gap: 4px; }

.bar {
  width: 36px;
  height: 4px;
  border-radius: 2px;
  background: #e5e7eb;
  transition: background 0.3s;
}

.bar.active.weak { background: #ef4444; }
.bar.active.medium { background: #f59e0b; }
.bar.active.strong { background: #10b981; }

.strength-label { font-size: 12px; color: #6b7280; }

/* 验证码 */
.captcha-wrapper { margin: 4px 0; }

/* 提交按钮 */
.submit-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.15s, box-shadow 0.15s, opacity 0.2s;
  box-shadow: 0 4px 20px rgba(99, 102, 241, 0.4);
  letter-spacing: 0.5px;
  margin-top: 4px;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 28px rgba(99, 102, 241, 0.5);
}

.submit-btn:active:not(:disabled) { transform: translateY(0); }
.submit-btn:disabled { opacity: 0.7; cursor: not-allowed; }

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255,255,255,0.4);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

/* 分割线 */
.divider {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #d1d5db;
  font-size: 13px;
}

.divider::before, .divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: #e5e7eb;
}

/* 社交登录 */
.social-login { display: flex; gap: 12px; }

.social-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 11px;
  border: 1.5px solid #e5e7eb;
  border-radius: 10px;
  background: white;
  font-size: 14px;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
  font-weight: 500;
}

.social-btn svg { width: 20px; height: 20px; }
.social-btn.wechat { color: #1aad19; }
.social-btn:hover { border-color: #6366f1; color: #4f46e5; box-shadow: 0 2px 8px rgba(0,0,0,0.08); }

/* 切换提示 */
.switch-hint {
  text-align: center;
  font-size: 13px;
  color: #9ca3af;
  margin-top: 20px;
}

.switch-link {
  color: #6366f1;
  font-weight: 600;
  text-decoration: none;
}
.switch-link:hover { text-decoration: underline; }

/* 过渡动画 */
.slide-fade-enter-active, .slide-fade-leave-active {
  transition: all 0.25s ease;
}
.slide-fade-enter-from { opacity: 0; transform: translateX(20px); }
.slide-fade-leave-to { opacity: 0; transform: translateX(-20px); }

/* ── 响应式断点 ──────────────────────────────── */

/* 平板横屏 / 小桌面 (900px ~ 1100px)：压缩左侧宽度 */
@media (max-width: 1100px) {
  .login-container { max-width: 820px; }
  .brand-panel { flex: 0 0 300px; padding: 40px 28px; }
  .brand-name { font-size: 30px; }
  .form-panel { padding: 40px 36px; }
}

/* 平板竖屏 (600px ~ 900px)：左侧收窄，隐藏 features */
@media (max-width: 900px) {
  .login-container { max-width: 680px; min-height: 520px; }
  .brand-panel { flex: 0 0 240px; padding: 36px 24px; }
  .brand-name { font-size: 26px; }
  .brand-slogan { font-size: 12px; margin-bottom: 24px; }
  .brand-features { display: none; }
  .brand-logo svg { width: 44px; height: 44px; }
  .form-panel { padding: 36px 28px; }
  .form-title { font-size: 22px; }
  .fields-row { grid-template-columns: 1fr; }
}

/* 手机横屏 / 大手机 (480px ~ 700px)：单列，品牌条 */
@media (max-width: 700px) {
  .login-page { padding: 0; align-items: stretch; }
  .login-container { flex-direction: column; border-radius: 0; min-height: 100dvh; max-width: 100%; box-shadow: none; }
  .brand-panel { flex: none; padding: 28px 24px 22px; min-height: auto; }
  .brand-logo svg { width: 36px; height: 36px; margin-bottom: 10px; }
  .brand-name { font-size: 24px; margin-bottom: 4px; }
  .brand-slogan { font-size: 12px; margin-bottom: 16px; }
  .brand-features { flex-direction: row; flex-wrap: wrap; gap: 10px; }
  .feature-item { font-size: 12px; }
  .feature-icon { width: 22px; height: 22px; font-size: 10px; }
  .form-panel { padding: 24px 20px 32px; flex: 1; }
  .form-title { font-size: 20px; }
  .form-subtitle { font-size: 13px; }
  .fields-row { grid-template-columns: 1fr; }
}

/* 小屏手机 (<480px)：最大化空间利用，减少间距 */
@media (max-width: 480px) {
  .brand-panel { padding: 20px 16px 16px; }
  .brand-logo svg { width: 30px; height: 30px; margin-bottom: 8px; }
  .brand-name { font-size: 20px; }
  .brand-features { gap: 8px; }
  .feature-item { font-size: 11px; gap: 6px; }
  .feature-icon { width: 20px; height: 20px; }
  .form-panel { padding: 20px 16px 28px; }
  .auth-form { gap: 12px; }
  .tab-switcher { width: 100%; }
  .input-wrapper input { padding: 11px 0; font-size: 13px; }
  .submit-btn { padding: 13px; font-size: 14px; }
  .form-header { margin-bottom: 20px; }
}
</style>
