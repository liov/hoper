import { httpclient } from '@/api/common'
import { LocaleResp } from '@gen/pb/common/common.service'
import { CommonResp } from '@hopeio/utils/types'

type LocalePack = {
  version: string
  locale: string
  messages: Record<string, string>
}

export const fetchLocalePack = async (locale: string): Promise<LocalePack | undefined> => {
  return await httpclient.get<LocalePack>('/api/locale', { query: { locale }, decode: LocaleResp })
}

