<route lang="json5">
{
  style: {
    navigationStyle: 'custom',
    navigationBarTitleText: 'wopan列表',
    onReachBottomDistance: 50,
  },
}
</route>
<template>
  <view class="list">
    <wd-navbar fixed left-text="返回" left-arrow right-icon="tune" placeholder>
      <template v-slot:title>
        <wd-button class="button-box" size="small" @click="toPDir">上一级</wd-button>
        <wd-button class="button-box" size="small" @click="deleteAll">删除全部</wd-button>
        <wd-button class="button-box" size="small" @click="deleteChecked">删除</wd-button>
      </template>
      <template v-slot:right>
        <wd-icon name="swap" size="18" @click="formData.waterfall = !formData.waterfall" />
      </template>
    </wd-navbar>
    <uni-list :class="{ 'uni-list--waterfall': formData.waterfall }">
      <wd-checkbox-group
        v-model="checkList"
        @change="checkboxChange"
        shape="square"
        style="display: contents"
        inline
      >
        <uni-list-item
          :border="!formData.waterfall"
          class="uni-list-item--waterfall"
          title="文件列表"
          v-for="(item, index) in curDir.subFiles"
          :key="item.file.id"
          clickable
          @click="onClick(item, index)"
        >
          <template v-slot:header v-if="item.file.type == 1">
            <view
              class="uni-thumb file-picture"
              :class="{ 'file-picture-column': formData.waterfall }"
            >
              <image
                :src="item.file.name.endsWith('.avif') ? item.file.previewUrl : item.file.thumbUrl"
                mode="aspectFill"
              ></image>
            </view>
          </template>

          <template v-slot:body>
            <view class="file">
              <view>
                <view class="uni-title">
                  <text class="uni-ellipsis-2">
                    {{
                      item.file.name.length < 25
                        ? item.file.name
                        : item.file.name.slice(0, 12) + '···' + item.file.name.slice(-12)
                    }}
                  </text>
                </view>
              </view>
              <view v-if="item.file.type == 1">
                <text class="uni-tag">{{ (item.file.size / 1024).toFixed(2) }}kb</text>
                <text class="uni-tag">{{ item.file.createTime.slice(0, 8) }}</text>
                <wd-checkbox :modelValue="index" shape="square" />
              </view>
              <view v-if="item.file.type == 1 && !formData.waterfall">
                <view class="uni-note ellipsis">
                  {{
                    item.file.fid.length < 40
                      ? item.file.fid
                      : item.file.fid.slice(0, 20) + '···' + item.file.fid.slice(-20)
                  }}
                </view>
                <view class="uni-note ellipsis">
                  {{
                    item.file.id.length < 40
                      ? item.file.id
                      : item.file.id.slice(0, 20) + '···' + item.file.id.slice(-20)
                  }}
                </view>
              </view>
            </view>
          </template>
        </uni-list-item>
      </wd-checkbox-group>
    </uni-list>
    <uni-load-more
      v-if="formData.loading || formData.status === 'no-more'"
      :status="formData.status"
    />
    <wd-backtop :scrollTop="scrollTop"></wd-backtop>
  </view>
</template>

<script lang="ts" setup>
import { useWopanStore } from '@/store/wopan'
import * as wopan from 'diamond/wopan'
import { storeToRefs } from 'pinia'
import { FileNode } from '@/model/wopan'
import { onPullDownRefresh, onReachBottom, onPageScroll } from '@dcloudio/uni-app'

defineOptions({
  name: 'WopanList',
})

const scrollTop = ref<number>(0)
onPageScroll((e) => {
  scrollTop.value = e.scrollTop
})
const wopanStore = useWopanStore()
if (wopanStore.$state.accessToken === '') {
  uni.navigateTo({
    url: '/pages/wopan/login',
  })
}
if (wopanStore.$state.psToken === '') {
  uni.navigateTo({
    url: '/pages/wopan/login?psToken=1',
  })
}
const formData = reactive({
  waterfall: false, // 布局方向切换
  status: 'loading', // 加载状态
  tipShow: false,
  loading: false,
})
if (wopanStore.$state.curDir.subFiles.length === 0) {
  loadFiles()
}
function loadFiles() {
  formData.status = 'loading'
  formData.loading = true
  wopanStore.FileList().then(() => {
    console.log('FileList', wopanStore.$state.curDir.subFiles)
    if (!wopanStore.$state.curDir.hasMore) {
      formData.status = 'no-more'
    } else {
      formData.status = 'more'
    }
    formData.loading = false
  })
}
const { curDir } = storeToRefs(wopanStore)
console.log('curDir', curDir)
onPullDownRefresh(() => {
  formData.status = 'more'
  formData.tipShow = true
  //
  formData.tipShow = false
  uni.stopPullDownRefresh()
})

function onClick(file: FileNode, index: number) {
  if (file.file.type === 0) {
    wopanStore.$state.curDir = file
    if (wopanStore.$state.curDir.subFiles.length === 0) {
      loadFiles()
    }
  }
}
const checkList = ref([])
function checkboxChange(e) {
  console.log(e)
}
function deleteChecked() {
  const dirList: string[] = []
  const fileList: string[] = []
  for (const i of checkList.value) {
    if (wopanStore.$state.curDir.subFiles[i].file.type === 1) {
      dirList.push(wopanStore.$state.curDir.subFiles[i].file.id)
    } else {
      fileList.push(wopanStore.$state.curDir.subFiles[i].file.id)
    }
  }
  wopan.DeleteFile(wopan.SpaceType.Private, dirList, fileList).then(() => {
    for (const i of checkList.value) {
      wopanStore.$state.curDir.subFiles.splice(i, 1)
    }
  })
}
function deleteAll() {
  console.log('deleteAll')
}
function toPDir() {
  wopanStore.$state.curDir = wopanStore.$state.curDir.parent
}
onReachBottom(() => {
  console.log('onReachBottom')
  loadFiles()
})

function handleClickLeft() {
  uni.navigateBack()
}
</script>

<style lang="scss" scoped>
@import '../../style/uni-ui.scss';

page {
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  background-color: #efeff4;
  min-height: 100%;
  height: auto;
}
.uni-navbar__header-container {
  display: flex;
  align-self: center;
}
.tips {
  color: #67c23a;
  font-size: 14px;
  line-height: 40px;
  text-align: center;
  background-color: #f0f9eb;
  height: 0;
  opacity: 0;
  transform: translateY(-100%);
  transition: all 0.3s;
}

.tips-ani {
  transform: translateY(0);
  height: 40px;
  opacity: 1;
}

.file {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.file-picture {
  width: 100px;
  height: 100px;
}

.file-picture-column {
  width: 100%;
  height: 100px;
}

.file-price {
  font-size: 12px;
  color: #ff5a5f;
}

.file-price-text {
  font-size: 16px;
}

.hot-tag {
  background: #ff5a5f;
  border: none;
  color: #fff;
}

.button-box {
  margin-left: 5px;
}

.uni-link {
  flex-shrink: 0;
}

.ellipsis {
  display: flex;
  overflow: hidden;
}

.uni-ellipsis-1 {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.uni-ellipsis-2 {
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

// 默认加入 scoped ，所以外面加一层提升权重
.list {
  /* #ifdef H5 || APP-VUE */
  ::v-deep
    /* #endif */
  .uni-section-header {
    padding: 5px 10px;
  }
  /* #ifdef H5 || APP-VUE */
  ::v-deep
    /* #endif */
  .uni-list-item__container {
    padding: 5px 10px;
  }
  .wd-checkbox.is-inline {
    display: contents;
  }
  .uni-list--waterfall {
    /* #ifndef H5 || APP-VUE */
    // 小程序 编译后会多一层标签，而其他平台没有，所以需要特殊处理一下
    ::v-deep(.uni-list) {
      /* #endif */
      display: flex;
      flex-direction: row;
      flex-wrap: wrap;
      padding: 5px;
      box-sizing: border-box;

      /* #ifdef H5 || APP-VUE */
      // h5 和 app-vue 使用深度选择器，因为默认使用了 scoped ，所以样式会无法穿透
      ::v-deep
      /* #endif */
      .uni-list-item--waterfall {
        width: 50%;
        box-sizing: border-box;

        .uni-list-item__container {
          flex-direction: column;
        }
      }

      /* #ifndef H5 || APP-VUE */
    }

    /* #endif */
  }
}
</style>
