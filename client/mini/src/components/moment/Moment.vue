<template>
  <div class="moment">
    <div class="auth">
      <image class="avatar" :src="staticDir + user.avatarUrl" />
      <span class="name">{{ user.name }}</span>
      <span class="time">{{ date2s(moment.createdAt) }}</span>
    </div>

    <div class="content" @click="detail">
      <nut-textarea v-model="moment.content"
                    rows="1"
                    :autosize="maxHeight ? { maxHeight } : true"
                    readonly/>
    </div>
    <view class="imgs" v-if="moment.images">
      <image
        style="width: 70px;height: 70px;background: #fff;"
        v-for="(img, idx) in images"
        :key="idx"
        :src="img"
        lazyLoad
        @click=preview(img)
        class="img"
      />
    </view>
  </div>
</template>

<script setup lang="ts">
import { jump } from "@/router/utils";
import { STATIC_DIR as staticDir } from "@/plugins/config";
import { date2s } from "@/plugins/utils/time";
import { reactive, ref } from "vue";
import Taro from "@tarojs/taro";

const props = defineProps<{
  moment: any;
  user: any;
  maxHeight?: number;
}>();

const moment = reactive(props.moment);
const route = Taro.getCurrentInstance().router;
const showPreview = ref(false);

const images = props.moment.images
  ?.split(",").map((image) => staticDir + image);

function preview(img:string) {
  Taro.previewImage({urls:images as string[],current:img})
}

function detail() {
  jump(route!.path, 1, props.moment);
}
</script>

<style scoped lang="scss">
.moment {
  $twelvepx: 20px;
  $avatar: 30px;
.auth{
  height: 30px;
  .avatar {
    left: 10px;
    position: absolute;
    flex-shrink: 0;
    width: $avatar;
    height: $avatar;
    border-radius: 40px;
    margin: 0 16px;
  }

  .name {
    left: 60px;
    position: absolute;
  }

  .time {
    position: absolute;
    right: $twelvepx;
  }
}


  .content {
    width: 100%;

    h3 {
      margin: 0;
      font-size: 18px;
      line-height: 20px;
    }

    .arrow {
      position: absolute;
      bottom: 16px;
      right: 0;
    }

    .van-multi-ellipsis--l3 {
      margin: 13px 0 0;
      font-size: 14px;
      line-height: 20px;
    }
  }

  .imgs {
    padding: 0 11px;
  }

  .img {
    margin: 5px 5px;
  }
}
</style>
