<template>
  <div class="moment">
    <div class="auth">
      <img class="avatar" :src="staticDir + user.avatarUrl" />
      <span class="name">{{ user.name }}</span>
      <span class="time">{{ dateFmtDateTime(moment.createdAt) }}</span>
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
    <lazy-component class="imgs" v-if="moment.images">
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
    <Action :content="moment" :type="1"></Action>
  </div>
</template>

<script setup lang="ts">
import { showImagePreview } from "vant";
import Action from "@/mixin/h5/components/action/Action.vue";
import { jump } from "@/mixin/router/utils";
import { STATIC_DIR as staticDir } from "@/mixin/plugin/config";
import { useRoute } from "vue-router";
import { dateFmtDateTime } from "@hopeio/utils/time";
import { reactive } from "vue";

const props = defineProps<{
  moment: any;
  user: any;
  maxHeight?: number;
}>();

const moment = reactive(props.moment);
const route = useRoute();

const images = props.moment.images
  ?.split(",")
  .map((image:string) => staticDir + image);

function preview(idx: number) {
  showImagePreview({
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
  text-align: left;

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
