const rbBase = import.meta.env.VITE_SERVER_BASEURL || '//api.hoper.xyz'

export type RbFileEntry = {
  id: string
  name: string
  size: number
  thumbHash?: string
  thumb_hash?: string
}

function rbUrl(path: string) {
  const base = rbBase.startsWith('//') ? `https:${rbBase}` : rbBase
  return `${base}${path}`
}

function uniGet<T>(path: string, query?: Record<string, string>) {
  return new Promise<T>((resolve, reject) => {
    uni.request({
      url: rbUrl(path),
      method: 'GET',
      data: query,
      success: (res) => resolve(res.data as T),
      fail: (err) => reject(err),
    })
  })
}

export async function rbHealth() {
  return uniGet<Record<string, string>>('/rb/health')
}

export async function rbListFiles(path: string) {
  const data = await uniGet<{ entries?: RbFileEntry[] }>('/rb/v1/list', { path })
  return data.entries || []
}

export function rbThumbUrl(path: string, maxEdge = 256, hash = '') {
  const q = new URLSearchParams({ path, max_edge: String(maxEdge) })
  if (hash) {
    q.set('hash', hash)
  }
  return `${rbUrl('/rb/v1/thumb')}?${q.toString()}`
}

export function rbSignalWsUrl() {
  const base = rbBase.startsWith('//') ? `https:${rbBase}` : rbBase
  const u = new URL(base)
  u.protocol = u.protocol === 'https:' ? 'wss:' : 'ws:'
  u.pathname = '/rb/signal'
  u.search = ''
  return u.toString()
}
