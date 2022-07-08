<template>
  <div class="moment">
    <div class="auth">
      <img class="avatar" :src="staticDir + user.avatarUrl" />
      <span class="name">{{ user.name }}</span>
      <span class="time">{{ $date2s(moment.createdAt) }}</span>
    </div>
    <div class="content" @click="detail">
      <van-field
        v-model="moment.content"
        rows="1"
        :autosize="maxHeight ? { maxHeight } : true"
        readonly
        type="textarea"
      >
        <template #extra>
          <div class="arrow">
            <van-icon name="arrow-down" />
          </div>
        </template>
      </van-field>
    </div>
    <lazy-component class="imgs" v-if:="moment.images">
      <van-image
        width="100"
        height="100"
        v-for="(img, idx) in images"
        :key="idx"
        :src="img"
        lazy-load
        class="img"
        @click="preview(idx)"
      />
    </lazy-component>
    <Action :content="moment" :type="1"> </Action>
  </div>
</template>

<script setup lang="ts">
import { ImagePreview } from "vant";
import Action from "@/components/action/Action.vue";
import { jump } from "@/router/utils";
import { STATIC_DIR as staticDir } from "@/plugin/config";
import { useRoute } from "vue-router";

const props = defineProps<{
  moment: any;
  user: any;
  maxHeight: number;
}>();

const route = useRoute();

const images = props.moment.images
  ?.split(",")
  .map((image) => staticDir + image);

function preview(idx: number) {
  ImagePreview({
    images,
    startPosition: idx,
    closeable: true,
  });
}
function detail() {
  jump(route.path, 1, props.moment);
}
</script>

<style scoped lang="less">
.moment {
  @20px: 20px;
  @avatar: 30px;
  .name {
    left: 60px;
    position: absolute;
  }

  .time {
    position: absolute;
    right: @20px;
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

  .avatar {
    flex-shrink: 0;
    width: @avatar;
    height: @avatar;
    border-radius: 40px;
    position: relative;
    margin: 0 16px;
  }
  .imgs {
    padding: 0 11px;
  }
  .img {
    margin: 5px 5px;
  }
}
</style>
