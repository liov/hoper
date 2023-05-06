<template>
    <view>
    <uni-list :border="false">
        <view v-for="(moment,idx) in momentList" :key="idx"
              :set=" users[idx] = userStore.getUser(moment.userId)">

            <uni-card margin="2px 5px" padding="0">
                <uni-list-item direction="column">
                    <!-- 自定义 header -->
                    <template v-slot:header>
                        <uni-list-chat :avatar-circle="true" :title="users[idx].name" :avatar="staticDir+users[idx].avatarUrl"
                                       note="来自iPhone15 Pro Max" :time="moment.createdAt.slice(0,10)+' '+ moment.createdAt.slice(11)"></uni-list-chat>
                    </template>
                    <!-- 自定义 body -->
                    <template v-slot:body>
                        <text class="slot-text">{{ moment.content }}</text>
                        <view v-if="moment.images" class="uni-title-sub uni-ellipsis-2" :let="images = moment.images.split(',')">
                            <view v-if="images.length == 1">
                                <image :src="staticDir+images[0]" style="width: 60%;max-height: 500px" mode="aspectFill" />
                            </view>
                            <view v-else-if="(images.length = 2) || (images.length = 4)">
                                <uni-grid :column="2" :show-border="false" :square="false">
                                    <uni-grid-item v-for="(img, index) in images" :index="index" :key="index">
                                        <view class="grid-item-box">
                                            <image :src="staticDir+img" :style="{width: '100%',height: 360/images.length+'px'}" mode="aspectFill" />
                                        </view>
                                    </uni-grid-item>
                                </uni-grid>
                            </view>
                            <view v-else-if="images.length < 10">
                            <uni-grid :column="3" :show-border="false" :square="false">
                                <uni-grid-item v-for="(img, index) in images" :index="index" :key="index">
                                    <view class="grid-item-box">
                                        <image :src="staticDir+img" style="width: 100%;height: 90px" mode="aspectFill" />
                                    </view>
                                </uni-grid-item>
                            </uni-grid>
                            </view>
                            <view v-else>
                                <swiper class="swiper" :indicator-dots="true">
                                    <swiper-item v-for="idx in Math.ceil(images.length/9)">
                                        <uni-grid :column="3" :show-border="false" :square="false">
                                            <uni-grid-item v-for="(img, index) in images.slice(9*(idx-1),9*idx)" :index="index" :key="index">
                                                <view class="grid-item-box">
                                                    <image :src="staticDir+img" style="width: 100%;height: 90px" mode="aspectFill" />
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
    <uni-load-more iconType="circle" :status="loadMoreStatus" @clickLoadMore="clickLoadMore"/>
    </view>
</template>

<script setup lang="ts">
import Actions from "@/components/action"
import MomentService from "@/service/moment";
import type {Ref} from "vue";
import {ref} from "vue";
import {userStore} from "@/store"
import {STATIC_DIR as staticDir} from "@/env/config";
import type {Moment, MomentList} from "@/model/moment";
import type {User} from "@/model/user";
import {onPullDownRefresh, onReachBottom} from "@dcloudio/uni-app";
import type {PageRequest} from "@/service/request";

const listReq:PageRequest = {
    PageNo: 1, PageSize: 10
}

const users: User[] = [];
const momentList: Ref<Moment[]> = ref([]);
const loadMoreStatus = ref('more');


getMomentList().then((res) => {
    momentList.value = momentList.value.concat(res.list)
});

async function getMomentList():Promise<MomentList>{
    return MomentService.getMomentList(listReq.PageNo, listReq.PageSize)
}


onPullDownRefresh(async ()=> {
    console.log('refresh');
    listReq.PageNo = 1;
    const res = await getMomentList()
    momentList.value = res.list;
    uni.stopPullDownRefresh();
})


onReachBottom(async ()=>{
    listReq.PageNo++;
    loadMoreStatus.value = 'loading';
    const res =  await getMomentList();
    if (res.list) {
        momentList.value = momentList.value.concat(res.list)
        if( res.list < listReq.PageSize)  {
            loadMoreStatus.value = 'no-more';
            if (res.list.length == 0) listReq.PageNo--;
        } else {
            loadMoreStatus.value = 'more';
        }
    }else {
        loadMoreStatus.value = 'no-more';
    }
})

function clickLoadMore(){
    if (loadMoreStatus.value == 'no-more') uni.showToast({title:"没有更多了",icon:"none"});
}


</script>

<style scoped lang="scss">
.slot-text {
  margin: 12px 0;
}

.grid-item-box{
    margin: 1px;
}

.swiper{
    height: 276px;
}

:deep(.uni-list-chat__container) {
  padding: 0;
}

:deep(.uni-list-item__container) {
  padding: 12px 15px 2px 15px;
}
</style>
