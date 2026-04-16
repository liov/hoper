import { httpclient } from '@/api/common'
import { CommonResp } from '@hopeio/utils/types'

type LocalePack = {
  version: string
  locale: string
  messages: Record<string, string>
}

export const fetchLocalePack = async (locale: string): Promise<LocalePack | undefined> => {
  const resp = await httpclient.get<CommonResp<LocalePack>>('/api/locale', { query: { locale } })
  return resp.data
}

