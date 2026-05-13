<template>
  <view class="page">
    <wd-navbar fixed placeholder title="远程配对" left-arrow @click-left="back" />
    <view class="form">
      <wd-input v-model="room" label="房间码" placeholder="与 Agent 一致" />
      <wd-input v-model="path" label="目录 path" placeholder="Agent 根目录或子路径" />
      <wd-button block type="primary" :loading="loading" @click="pair">配对信令</wd-button>
      <wd-button block plain :disabled="!paired" @click="goGallery">进入相册</wd-button>
      <text class="hint">信令 {{ signalUrl }}</text>
      <text v-if="relay" class="hint">中继 {{ relay }}</text>
      <text v-if="status" class="hint">{{ status }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { onUnload, ref } from 'vue'
import { rbSignalWsUrl } from '@/api/remotebrowse'
import { encodeRegisterViewer, parseSignalEnvelope } from '@/api/rb_signal'

const room = ref('demo')
const path = ref('')
const loading = ref(false)
const paired = ref(false)
const relay = ref('')
const status = ref('')
const signalUrl = rbSignalWsUrl()
let socket: UniApp.SocketTask | null = null

function pair() {
  loading.value = true
  status.value = '连接信令…'
  socket?.close({})
  socket = uni.connectSocket({ url: signalUrl })
  socket.onOpen(() => {
    socket?.send({ data: encodeRegisterViewer(room.value) })
    status.value = '已注册 viewer，等待 Agent…'
  })
  socket.onMessage((msg) => {
    const data = msg.data as ArrayBuffer
    const env = parseSignalEnvelope(data)
    if (env.error) {
      status.value = env.error
      loading.value = false
      return
    }
    if (env.relay) {
      paired.value = true
      relay.value = `${env.relay.relayHost}:${env.relay.relayPort}`
      status.value = '中继已下发，可进入相册（列表走 HTTP）'
      loading.value = false
    }
  })
  socket.onError(() => {
    status.value = '信令连接失败'
    loading.value = false
  })
}

function goGallery() {
  uni.navigateTo({
    url: `/pages/remotebrowse/gallery?room=${encodeURIComponent(room.value)}&path=${encodeURIComponent(path.value)}`,
  })
}

function back() {
  uni.navigateBack()
}

onUnload(() => {
  socket?.close({})
})
</script>

<style scoped>
.page {
  padding: 12px;
}
.form {
  margin-top: 48px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.hint {
  font-size: 12px;
  color: #666;
  word-break: break-all;
}
</style>
