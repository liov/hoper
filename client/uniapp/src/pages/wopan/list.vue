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
    <wd-navbar
      fixed
      left-text="返回"
      left-arrow
      @click-left="handleClickLeft"
      right-icon="tune"
      placeholder
    >
      <template v-slot:title>
        <text>文件列表</text>
      </template>
      <template v-slot:right>
        <wd-button @click="deleteAll" custom-class="title-right" type="primary" size="small" round>
          删除全部
        </wd-button>
        <wd-icon
          name="swap"
          size="18"
          @click="formData.waterfall = !formData.waterfall"
          custom-class="title-right"
        />
      </template>
    </wd-navbar>
    <uni-breadcrumb separator=">">
      <uni-breadcrumb-item
        v-for="(file, index) in breadcrumb"
        :key="index"
        @click="breadcrumbTo(index)"
      >
        {{ file.file.name }}
      </uni-breadcrumb-item>
    </uni-breadcrumb>
    <uni-list :class="{ 'uni-list--waterfall': formData.waterfall }">
      <wd-checkbox-group
        v-model="checkList"
        shape="square"
        style="display: contents"
        inline
      >
        <wd-swipe-action  v-for="(item, index) in wopanStore.$state.curDir.subFiles" :key="item.file.id"  custom-class="uni-list-item--waterfall">
        <uni-list-item
          :border="!formData.waterfall"
          title="文件列表"
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
                  <text class="uni-ellipsis-2" :style="{ color: item.read ? '#6FD8D6' : '#303030' }">
                    {{
                      formData.waterfall?
                      item.file.name.length < 18 ?
                        item.file.name : item.file.name.slice(0, 13) + '···'+ item.file.name.slice(-4):
                        item.file.name.length < 30 ?
                          item.file.name : item.file.name.slice(0, 24) + '···' + item.file.name.slice(-4)
                    }}
                  </text>
                  <wd-checkbox :modelValue="index" shape="square" @click.stop/>
                </view>
              </view>
              <view v-show="item.file.type == 1">
                <text class="uni-tag">{{ (item.file.size / 1024).toFixed(2) }}kb</text>
                <text class="uni-tag">{{ item.file.createTime }}</text>
              </view>
              <view v-show="item.file.type == 1 && !formData.waterfall">
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
          <template #right>
            <wd-button size="small" @click="deleteFile(index)">删除</wd-button>
          </template>
        </wd-swipe-action>
      </wd-checkbox-group>
    </uni-list>
    <uni-load-more
      v-show="formData.loading || formData.status === 'no-more'"
      :status="formData.status"
    />
    <wd-backtop :scrollTop="scrollTop"></wd-backtop>
<!--    <wd-fab  draggable type="primary" position="left-bottom">
      <wd-button
        @click="deleteChecked"
        custom-class="custom-button"
        type="success"
        round
        size="small"
      >
        删除选中
      </wd-button>
    </wd-fab>-->
    <wd-fab
      v-show="checkList.length>0"
      type="primary"
      position="right-bottom"
      :expandable="false"
      @click="deleteChecked"
      inactiveIcon="delete"
    ></wd-fab>
  </view>
</template>

<script lang="ts" setup>
import { useWopanStore } from '@/store/wopan'
import * as wopan from '@hopeio/utils/wopan'
import { storeToRefs } from 'pinia'
import { FileNode } from '@/model/wopan'
import { onPullDownRefresh, onReachBottom, onPageScroll } from '@dcloudio/uni-app'
import { useMessage } from 'wot-design-uni'
const message = useMessage()
defineOptions({
  name: 'WopanList',
})

const scrollTop = ref<number>(0)
onPageScroll((e) => {
  scrollTop.value = e.scrollTop
})
const wopanStore = useWopanStore()
wopanStore.checkToken()
const formData = reactive({
  waterfall: false, // 布局方向切换
  status: 'loading', // 加载状态
  tipShow: false,
  loading: false,
})
const breadcrumb = ref([wopanStore.$state.file])

if (wopanStore.$state.curDir.subFiles.length === 0) {
  loadFiles()
}
function loadFiles() {
  formData.status = 'loading'
  formData.loading = true
  wopanStore.FileList().then(() => {
    if (!wopanStore.$state.curDir.hasMore) {
      formData.status = 'no-more'
    } else {
      formData.status = 'more'
    }
    formData.loading = false
  })
}

function breadcrumbTo(index: number) {
  wopanStore.$state.curDir = breadcrumb.value[index]
  breadcrumb.value = breadcrumb.value.slice(0, index + 1)
}

onPullDownRefresh(() => {
  formData.status = 'more'
  formData.tipShow = true
  //
  formData.tipShow = false
  uni.stopPullDownRefresh()
})

function onClick(file: FileNode, index: number) {
  file.read = true
  if (file.file.type === 0) {
    wopanStore.$state.curDir = file
    wopanStore.$state.viewList.push(index)
    if (wopanStore.$state.curDir.subFiles.length === 0) {
      loadFiles()
    }
    breadcrumb.value.push(wopanStore.$state.curDir)
    checkList.value = []
  }else{
    uni.navigateTo({
      url: '/pages/wopan/view?index=' + index,
    })
  }
}
const checkList = ref([])

function deleteChecked() {
  console.log('delete checked')
  wopanStore.deleteFiles(checkList.value)
}
function deleteAll() {
  message
    .confirm({
      msg: '确认删除',
      title: '删除全部',
    })
    .then(() => {
      console.log('点击了确定按钮')
    })
    .catch(() => {
      console.log('点击了取消按钮')
    })
}
function toPDir() {
  wopanStore.$state.curDir = wopanStore.$state.curDir.parent
  breadcrumb.value.pop()
  checkList.value = []
  wopanStore.$state.viewList.pop()
}
onReachBottom(() => {
  console.log('onReachBottom')
  loadFiles()
})

function handleClickLeft() {
  if (wopanStore.$state.curDir.file.id !== '0') {
    toPDir()
  } else {
    uni.navigateBack()
  }
}
function deleteFile(index:number) {
  wopanStore.deleteCurDirFile(index);
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
.title-right {
  margin-left: 10px;
}
/* #ifdef H5 || APP-VUE */
::v-deep
  /* #endif */
.uni-breadcrumb-item--slot {
  padding: 0;
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
/* #ifdef H5 || APP-VUE */
::v-deep
  /* #endif */
.wd-checkbox__shape{
  position: absolute;
  right: 0;
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
  /* #ifdef H5 || APP-VUE */
  ::v-deep
    /* #endif */
  .wd-swipe-action__right{
    display: flex;
    justify-content: center;
    align-items: center;
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
