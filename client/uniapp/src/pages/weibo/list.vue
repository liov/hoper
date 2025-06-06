<route lang="json5">
{
style: {
navigationStyle: 'custom',
navigationBarTitleText: 'weibo列表',
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
        <wd-button  custom-class="title-right" type="primary" size="small" round>
          全选
        </wd-button>
        <wd-icon
          name="swap"
          size="18"
          @click="formData.waterfall = !formData.waterfall"
          custom-class="title-right"
        />
      </template>
    </wd-navbar>

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
import { onPullDownRefresh, onReachBottom, onPageScroll } from '@dcloudio/uni-app'
import { useMessage } from 'wot-design-uni'
import WeiboService from "@/api/weibo";
const message = useMessage()
defineOptions({
  name: 'WeiboList',
})

const scrollTop = ref<number>(0)
onPageScroll((e) => {
  scrollTop.value = e.scrollTop
})

onLoad(() => {
  formData.uid = uni.getStorageSync('uid')
  getList()
})

const formData = reactive({
  waterfall: false, // 布局方向切换
  status: 'loading', // 加载状态
  tipShow: false,
  loading: false,
  uid: 0,
  page:1,
})

async function getList() {
  formData.status = 'loading'
  formData.loading = true
  const res= await WeiboService.list(formData.uid,formData.page);
 console.log(res);
  formData.loading = false
}

onPullDownRefresh(() => {
  formData.status = 'more'
  formData.tipShow = true
  //
  formData.tipShow = false
  uni.stopPullDownRefresh()
})


const checkList = ref([])

function deleteChecked() {
  console.log('delete checked')
}


onReachBottom(() => {
  console.log('onReachBottom')
  getList()
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
