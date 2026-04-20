import { createI18n } from 'vue-i18n'
import { fetchLocalePack } from '../api/i18n'

import en from './en.json'
import zhHans from './zh-Hans.json' // 简体中文

type LocaleMessages = Record<string, string>
type LocalePack = { locale: string; messages: LocaleMessages }

const i18n = createI18n({
  legacy: false, // 启用 Composition API 模式，useI18n() 需要此配置
  locale: uni.getLocale(), // 获取已设置的语言，fallback 语言需要再 manifest.config.ts 中设置
  fallbackLocale: 'zh-Hans',
  flatJson: true,
  messages: {
    en,
    'zh-Hans': zhHans, // key 不能乱写，查看截图 screenshots/i18n.png
  },
})

const getCacheKey = (locale: string) => `i18n:${locale}`

const readCachedLocalePack = (locale: string): LocalePack | undefined => {
  const cache = uni.getStorageSync(getCacheKey(locale))
  if (!cache || typeof cache !== 'object') return undefined
  return cache as LocalePack
}

export const syncLocaleMessages = async (inputLocale?: string) => {
  const locale = inputLocale || uni.getLocale()
  const cached = readCachedLocalePack(locale)
  i18n.global.mergeLocaleMessage(locale, cached?.messages)
  try {
    const remote = await fetchLocalePack(locale).catch((e) => console.log(e))
    if (!remote) return
    i18n.global.mergeLocaleMessage(locale, remote.messages)
    uni.setStorageSync(getCacheKey(locale), remote)
  } catch (error) {
    console.error('[i18n] syncLocaleMessages error', error)
  }
}

export const setLocaleAndSync = async (locale: string) => {
  uni.setLocale(locale)
  i18n.global.locale.value = locale
  await syncLocaleMessages(locale)
}

/**
 * 非 vue 文件使用这个方法
 * @param { string } localeKey 多语言的key，eg: "app.name"
 */
export const translate = (localeKey: string, params?: Record<string, any>) => {
  if (!localeKey) {
    console.error(`[i18n] Function translate(), localeKey param is required`)
    return ''
  }
  if (i18n.global.te(localeKey)) return i18n.global.t(localeKey, params) as string
  return localeKey
}

export function tarbarI18n() {
  uni.setTabBarItem({
    index: 0,
    text: translate('tabbar.home'),
  });
  uni.setTabBarItem({
    index: 1,
    text: translate('tabbar.record'),
  });
}


export default i18n
