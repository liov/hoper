import { createI18n } from 'vue-i18n'
import { fetchLocalePack } from '../api/i18n'

import en from './en.json'
import zhHans from './zh-Hans.json' // 简体中文

type Locale = 'en' | 'zh-Hans'
type LocaleMessages = Record<string, string>
type LocalePack = { locale: string; messages: LocaleMessages }

const fallbackMessages: Record<Locale, LocaleMessages> = {
  en,
  'zh-Hans': zhHans, // key 不能乱写，查看截图 screenshots/i18n.png
}

const normalizeLocale = (locale?: string): Locale => {
  if (!locale) return 'zh-Hans'
  if (locale === 'zh' || locale.startsWith('zh-')) return 'zh-Hans'
  if (locale.startsWith('en')) return 'en'
  return locale === 'en' ? 'en' : 'zh-Hans'
}

const currentLocale = normalizeLocale(uni.getLocale())

const i18n = createI18n({
  legacy: false, // 启用 Composition API 模式，useI18n() 需要此配置
  locale: currentLocale, // 获取已设置的语言，fallback 语言需要再 manifest.config.ts 中设置
  fallbackLocale: 'zh-Hans',
  messages: fallbackMessages,
})

const getCacheKey = (locale: Locale) => `i18n:${locale}`

const mergeLocaleMessages = (locale: Locale, messages?: LocaleMessages) => {
  if (!messages || Object.keys(messages).length === 0) return
  i18n.global.mergeLocaleMessage(locale, messages)
}

const readCachedLocalePack = (locale: Locale): LocalePack | undefined => {
  const cache = uni.getStorageSync(getCacheKey(locale))
  if (!cache || typeof cache !== 'object') return undefined
  return cache as LocalePack
}

export const syncLocaleMessages = async (inputLocale?: string) => {
  const locale = normalizeLocale(inputLocale || uni.getLocale())
  const cached = readCachedLocalePack(locale)
  mergeLocaleMessages(locale, cached?.messages)
  try {
    const remote = await fetchLocalePack(locale)
    if (!remote) return
    mergeLocaleMessages(locale, remote.messages)
    uni.setStorageSync(getCacheKey(locale), remote)
  } catch (error) {
    console.error('[i18n] syncLocaleMessages error', error)
  }
}

export const setLocaleAndSync = async (locale: string) => {
  const normalized = normalizeLocale(locale)
  uni.setLocale(normalized)
  i18n.global.locale.value = normalized
  await syncLocaleMessages(normalized)
}

/**
 * 非 vue 文件使用这个方法
 * @param { string } localeKey 多语言的key，eg: "app.name"
 */
export const translate = (localeKey: string) => {
  if (!localeKey) {
    console.error(`[i18n] Function translate(), localeKey param is required`)
    return ''
  }
  if (i18n.global.te(localeKey)) return i18n.global.t(localeKey) as string
  return localeKey
}

/**
 * formatString('已阅读并同意{0}和{1}','用户协议','隐私政策') -> 已阅读并同意用户协议和隐私政策
 * @param template
 * @param values
 * @returns
 */
export function formatString(template: string, ...values: any) {
  console.log(template, values)
  // 使用map来替换{0}, {1}, ...等占位符
  return template.replace(/{(\d+)}/g, (match, index) => {
    const value = values[index]
    return value !== undefined ? value : match
  })
}

/**
 * formatI18n('我是{name},身高{detail.height},体重{detail.weight}',{name:'张三',detail:{height:178,weight:'75kg'}})
 * 暂不支持数组
 * @param template 多语言模板字符串，eg: `我是{name}`
 * @param obj 需要传递的数据对象，里面的key与多语言字符串对应，eg: `{name:'菲鸽'}`
 * @returns
 */
export function formatI18n(template: string, obj: Record<string, any>) {
  const match = /\{(.*?)\}/g.exec(template)
  if (match) {
    const variableList = match[0].replace('{', '').replace('}', '').split('.')
    let result: any = obj
    for (let i = 0; i < variableList.length; i++) {
      result = result[variableList[i]] || ''
    }
    return formatI18n(template.replace(match[0], String(result)), obj)
  } else {
    return template
  }
}

export default i18n
