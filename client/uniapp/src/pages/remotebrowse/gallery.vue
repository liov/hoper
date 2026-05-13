<template>
  <view class="page">
    <wd-navbar fixed placeholder :title="`相册 · ${room}`" left-arrow @click-left="back" />
    <view class="status">{{ status }}</view>
    <view v-if="preview" class="preview">
      <image :src="preview" mode="aspectFit" class="preview-img" />
      <text class="preview-name">{{ previewName }}</text>
    </view>
    <view class="grid">
      <view v-for="(item, idx) in entries" :key="item.id || item.name" class="cell" @click="select(idx)">
        <image :src="thumbOf(item)" mode="aspectFill" class="thumb" />
        <text class="name">{{ item.name }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { rbListFiles, rbThumbUrl, type RbFileEntry } from '@/api/remotebrowse'

const room = ref('demo')
const path = ref('')
const entries = ref<RbFileEntry[]>([])
const status = ref('')
const preview = ref('')
const previewName = ref('')

onLoad((q) => {
  room.value = decodeURIComponent(String(q?.room || 'demo'))
  path.value = decodeURIComponent(String(q?.path || ''))
  load()
})

async function load() {
  status.value = '加载中…'
  try {
    entries.value = await rbListFiles(path.value)
    status.value = `共 ${entries.value.length} 项`
    if (entries.value.length) {
      select(0)
    }
  } catch (e) {
    status.value = `失败 ${e}`
  }
}

function thumbOf(item: RbFileEntry) {
  const p = item.id || item.name
  const hash = item.thumbHash || item.thumb_hash || ''
  return rbThumbUrl(p, 256, hash)
}

function select(idx: number) {
  const item = entries.value[idx]
  if (!item) {
    return
  }
  previewName.value = item.name
  preview.value = thumbOf(item)
}

function back() {
  uni.navigateBack()
}
</script>

<style scoped>
.page {
  padding: 8px;
}
.status {
  margin-top: 44px;
  font-size: 12px;
  color: #666;
}
.preview {
  margin: 8px 0;
  height: 220px;
  background: #111;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.preview-img {
  flex: 1;
  width: 100%;
}
.preview-name {
  color: #fff;
  font-size: 12px;
  padding: 6px;
}
.grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 6px;
}
.cell {
  background: #f5f5f5;
  border-radius: 6px;
  overflow: hidden;
}
.thumb {
  width: 100%;
  height: 96px;
  display: block;
}
.name {
  font-size: 11px;
  padding: 4px;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
