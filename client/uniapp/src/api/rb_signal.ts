function writeVarint(n: number): number[] {
  const out: number[] = []
  let v = n >>> 0
  while (v >= 0x80) {
    out.push((v & 0x7f) | 0x80)
    v >>>= 7
  }
  out.push(v)
  return out
}

function writeStringField(field: number, s: string): number[] {
  const body = new TextEncoder().encode(s)
  const tag = writeVarint((field << 3) | 2)
  const len = writeVarint(body.length)
  return [...tag, ...len, ...body]
}

function writeMessageField(field: number, body: number[]): number[] {
  const tag = writeVarint((field << 3) | 2)
  const len = writeVarint(body.length)
  return [...tag, ...len, ...body]
}

export function encodeRegisterViewer(room: string): ArrayBuffer {
  const inner = [...writeStringField(1, room), ...writeStringField(2, 'viewer')]
  const env = writeMessageField(1, inner)
  return new Uint8Array(env).buffer
}

export type RelayTokenInfo = { sessionId: string; relayHost: string; relayPort: number }

function readVarint(buf: Uint8Array, off: number): [number, number] {
  let shift = 0
  let n = 0
  let i = off
  while (i < buf.length) {
    const b = buf[i++]
    n |= (b & 0x7f) << shift
    if ((b & 0x80) === 0) {
      return [n, i]
    }
    shift += 7
  }
  throw new Error('varint eof')
}

function readString(buf: Uint8Array, off: number): [string, number] {
  const [len, p] = readVarint(buf, off)
  const end = p + len
  return [new TextDecoder().decode(buf.subarray(p, end)), end]
}

function parseRelayToken(buf: Uint8Array): RelayTokenInfo {
  let sessionId = ''
  let relayHost = ''
  let relayPort = 0
  let i = 0
  while (i < buf.length) {
    const [tag, p] = readVarint(buf, i)
    i = p
    const field = tag >>> 3
    const wire = tag & 7
    if (wire === 2) {
      const [s, n] = readString(buf, i)
      i = n
      if (field === 1) sessionId = s
      if (field === 2) relayHost = s
    } else if (wire === 0) {
      const [v, n] = readVarint(buf, i)
      i = n
      if (field === 3) relayPort = v
    } else {
      break
    }
  }
  return { sessionId, relayHost, relayPort }
}

export function parseSignalEnvelope(buf: ArrayBuffer): { relay?: RelayTokenInfo; error?: string } {
  const u8 = new Uint8Array(buf)
  let i = 0
  while (i < u8.length) {
    const [tag, p] = readVarint(u8, i)
    i = p
    const field = tag >>> 3
    const wire = tag & 7
    if (wire !== 2) {
      break
    }
    const [len, p2] = readVarint(u8, i)
    i = p2
    const slice = u8.subarray(i, i + len)
    i += len
    if (field === 6) {
      return { relay: parseRelayToken(slice) }
    }
    if (field === 7) {
      return { error: new TextDecoder().decode(slice) }
    }
  }
  return {}
}
