<template>
  <div class="comment">
    <div class="auth">
      <img class="avatar" :src="staticDir + user.avatarUrl" />
      <span class="name">{{ user.name }}</span>
      <span class="time">{{ date2s(comment.createdAt) }}</span>
      <div class="like">
        <van-icon
          :name="comment.likeId > 0 ? 'like' : 'like-o'"
          :color="comment.likeId > 0 ? '#D91E46' : ''"
          @click="like"
        /><span class="count" v-if="comment.ext">{{
          comment.ext.likeCount
        }}</span>
      </div>
    </div>
    <div class="content">
      <van-field
        v-model="comment.content"
        rows="1"
        :autosize="{ maxHeight: 200 }"
        readonly
        type="textarea"
        @touchstart="onTouchStart"
        @touchend="onTouchEnd"
      >
      </van-field>
    </div>
    <lazy-component class="imgs" v-if:="comment.image">
      <van-image
        width="100"
        height="100"
        :src="staticDir + comment.image"
        lazy-load
        class="img"
        @click="preview"
      />
    </lazy-component>
  </div>
</template>

<script setup lang="ts">
import { showImagePreview } from "vant";
import { date2s } from "diamond/time";
import axios from "axios";
import emitter from "@/mixin/plugin/emitter";
import { STATIC_DIR as staticDir } from "@/mixin/plugin/config";
import { reactive } from "vue";

const props = defineProps<{
  comment: any;
  user: any;
}>();

const comment = reactive(props.comment);
const images = reactive([]);
let startTime = 0;
let stopTime = 0;
let timer;
const timeout = 500;

function preview() {
  showImagePreview({
    images: [comment.image],
    startPosition: 0,
    closeable: true,
  });
}
async function like() {
  const api = `/api/v1/action/like`;
  const id = comment.id;
  const likeId = comment.likeId;
  if (likeId > 0) {
    await axios.delete(`${api}/${likeId}`);
    comment.likeId = 0;
  } else {
    const res = await axios.post(api, {
      refId: id,
      type: 7,
      action: 2,
    });
    comment.likeId = res.data.data.id;
  }
}
function onComment() {
  emitter.emit("onComment", {
    replyId: comment.id,
    rootId: comment.rootId,
    recvId: comment.userId,
  });
}
function onTouchStart() {
  startTime = Date.now();
  timer = setTimeout(longPress, timeout);
}
function onTouchEnd(e: Event) {
  stopTime = Date.now();
  clearTimeout(timer);
  if (stopTime - startTime < timeout) {
    onComment();
    e.preventDefault();
  }
  return false;
}
function longPress() {
  emitter.emit("more-show", { type: 7, refId: comment.id });
}
</script>

<style scoped lang="less">
.comment {
  @20px: 20px;
  @avatar: 30px;
  .name {
    left: 60px;
    position: absolute;
  }

  .time {
    position: absolute;
    right: 80px;
  }
  .like {
    position: absolute;
    right: 20px;
    top: 0;
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
