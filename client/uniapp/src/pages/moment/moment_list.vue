<route lang="json5">
{
  style: {
    navigationBarTitleText: '瞬间',
  },
}
</route>
<template>
  <view>
    <uni-list :border="false">
      <view v-for="(moment, midx) in momentList" :key="midx">
        <uni-card margin="2px 5px" padding="0">
          <uni-list-item direction="column">
            <!-- 自定义 header -->
            <template v-slot:header>
              <uni-list-chat
                :avatar-circle="true"
                :title="moment.user.name"
                :avatar="staticDir + moment.user.avatarUrl"
                note="来自iPhone15 Pro Max"
                :time="moment.createdAt.slice(0, 10) + ' ' + moment.createdAt.slice(11, 19)"
              ></uni-list-chat>
            </template>
            <!-- 自定义 body -->
            <template v-slot:body>
              <text class="slot-text">{{ moment.content }}</text>
              <view v-if="moment.images" class="uni-title-sub uni-ellipsis-2">
                <view v-if="moment.imagesUrls.length == 1">
                  <image
                    :src="staticDir + moment.imagesUrls[0]"
                    style="width: 60%; max-height: 500px"
                    mode="aspectFill"
                  />
                </view>
                <view v-else-if="(moment.imagesUrls.length = 2) || (moment.imagesUrls.length = 4)">
                  <uni-grid :column="2" :show-border="false" :square="false">
                    <uni-grid-item
                      v-for="(img, iidx) in moment.imagesUrls"
                      :index="iidx"
                      :key="iidx"
                    >
                      <view class="grid-item-box">
                        <image
                          :src="staticDir + img"
                          :style="{ width: '100%', height: 360 / moment.imagesUrls.length + 'px' }"
                          mode="aspectFill"
                          @click="preview(midx, iidx)"
                        />
                      </view>
                    </uni-grid-item>
                  </uni-grid>
                </view>
                <view v-else-if="moment.imagesUrls.length < 10">
                  <uni-grid :column="3" :show-border="false" :square="false">
                    <uni-grid-item
                      v-for="(img, iidx) in moment.imagesUrls"
                      :index="iidx"
                      :key="iidx"
                    >
                      <view class="grid-item-box">
                        <image
                          :src="staticDir + img"
                          style="width: 100%; height: 90px"
                          mode="aspectFill"
                        />
                      </view>
                    </uni-grid-item>
                  </uni-grid>
                </view>
                <view v-else>
                  <swiper class="swiper" :indicator-dots="true">
                    <swiper-item v-for="sidx in Math.ceil(moment.imagesUrls.length / 9)">
                      <uni-grid :column="3" :show-border="false" :square="false">
                        <uni-grid-item
                          v-for="(img, iidx) in moment.imagesUrls.slice(9 * (sidx - 1), 9 * sidx)"
                          :index="iidx"
                          :key="iidx"
                        >
                          <view class="grid-item-box">
                            <image
                              :src="staticDir + img"
                              style="width: 100%; height: 90px"
                              mode="aspectFill"
                            />
                          </view>
                        </uni-grid-item>
                      </uni-grid>
                    </swiper-item>
                  </swiper>
                </view>
              </view>
            </template>
            <!-- 自定义 footer-->
          </uni-list-item>
          <template v-slot:actions>
            <Actions></Actions>
          </template>
        </uni-card>
      </view>
    </uni-list>
    <uni-load-more iconType="circle" :status="loadMoreStatus" @clickLoadMore="clickLoadMore" />
    <view class="popup-content">
      <!-- 普通弹窗 -->
      <uni-popup ref="popup" background-color="#fff">
        <view class="popup-content">
          <uni-swiper-dot
            class="uni-swiper-dot-box"
            :info="previewImgs"
            :current="current"
            field="content"
            mode="round"
          >
            <swiper style="width: 100%; height: 100%" @change="change">
              <swiper-item v-for="(img, iidx) in previewImgs" :key="iidx">
                <view class="swiper-item">
                  <image
                    :src="staticDir + img"
                    mode="aspectFit"
                    style="width: 100vw; height: 100vh"
                  />
                </view>
              </swiper-item>
            </swiper>
          </uni-swiper-dot>
        </view>
      </uni-popup>
    </view>
  </view>
</template>

<script setup lang="ts">
import Actions from '@/components/action.vue'
import MomentService from '@/service/moment'
import { userStore } from '@/store'
import { STATIC_DIR as staticDir } from '@/env/config'
import type { Moment, MomentList } from '@/model/moment'
import type { User } from '@/model/user'
import { onPullDownRefresh, onReachBottom } from '@dcloudio/uni-app'
import type { PageRequest } from '@/service/request'

const listReq: PageRequest = {
  PageNo: 1,
  PageSize: 10,
}

const momentList: Ref<Moment[]> = ref([])
const loadMoreStatus = ref('more')

getMomentList().then((res) => {
  momentList.value = momentList.value.concat(res.list)
})

async function getMomentList(): Promise<MomentList> {
  return MomentService.getMomentList(listReq.PageNo, listReq.PageSize)
}

onPullDownRefresh(async () => {
  console.log('refresh')
  listReq.PageNo = 1
  const res = await getMomentList()
  momentList.value = res.list
  uni.stopPullDownRefresh()
})

onReachBottom(async () => {
  if (loadMoreStatus.value === 'no-more') return
  listReq.PageNo++
  loadMoreStatus.value = 'loading'
  const res = await getMomentList()
  if (res.list) {
    momentList.value = momentList.value.concat(res.list)
    if (res.list.length < listReq.PageSize) {
      loadMoreStatus.value = 'no-more'
      if (res.list.length === 0) listReq.PageNo--
    } else {
      loadMoreStatus.value = 'more'
    }
  } else {
    loadMoreStatus.value = 'no-more'
    listReq.PageNo--
  }
})

function clickLoadMore() {
  if (loadMoreStatus.value === 'no-more') uni.showToast({ title: '没有更多了', icon: 'none' })
}

const popup = ref(null)
const previewImgs: Ref<string[]> = ref([])
const current = ref(0)
function preview(midx: number, iidx: number) {
  console.log(midx, iidx)
  previewImgs.value = momentList.value[midx].imagesUrls
  console.log(previewImgs.value)
  current.value = iidx
  popup.value.open('center')
}
function change(e) {
  current.value = e.detail.current
}
</script>

<style scoped lang="scss">
.slot-text {
  margin: 12px 0;
}

.grid-item-box {
  margin: 1px;
}

.swiper {
  height: 276px;
}

:deep(.uni-list-chat__container) {
  padding: 0;
}

:deep(.uni-list-item__container) {
  padding: 12px 15px 2px 15px;
}

@mixin flex {
  /* #ifndef APP-NVUE */
  display: flex;
  /* #endif */
  flex-direction: row;
}

@mixin height {
  /* #ifndef APP-NVUE */
  height: 100%;
  /* #endif */
  /* #ifdef APP-NVUE */
  flex: 1;
  /* #endif */
}

.popup-content {
  width: 100vw;
  height: 100vh;
}

.uni-swiper-dot-box {
  width: 100%;
  height: 100%;
  margin: 0 auto;
  margin-top: 8px;
}
/* 当屏幕宽度在 500px 以下时采用这种布局 */
@media screen and (max-width: 500px) {
  .image {
    width: 100%;
  }
}
/* 当屏幕宽度在 900px 以上时采用这种布局 */
@media screen and (min-width: 900px) {
  /* CSS 属性 */
}

.swiper-item {
  text-align: center;
}
</style>
