<template>
  <div class="comment">
    <div class="auth">
      <img class="avatar" :src="user.avatarUrl" />
      <span class="name">{{ user.name }}</span>
      <span class="time">{{ $date2s(comment.createdAt) }}</span>
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
        :src="comment.image"
        lazy-load
        class="img"
        @click="preview"
      />
    </lazy-component>
  </div>
</template>

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
import { ImagePreview } from "vant";
import Action from "@/components/action/Action.vue";
import axios from "axios";
import emitter from "@/plugin/emitter";
class Props {
  comment = prop<any>({ default: {} });
  user = prop<any>({});
}
@Options({ components: { Action } })
export default class Comment extends Vue.with(Props) {
  images = [];
  timeOutEvent = 0;
  preview() {
    ImagePreview({
      images: [this.comment.image],
      startPosition: 0,
      closeable: true,
    });
  }
  async like() {
    const api = `/api/v1/action/like`;
    const id = this.comment.id;
    const likeId = this.comment.likeId;
    if (likeId > 0) {
      await axios.delete(api, { data: { id: likeId } });
      this.comment.likeId = 0;
    } else {
      const res = await axios.post(api, {
        refId: id,
        type: 7,
        action: 2,
      });
      this.comment.likeId = res.data.details.id;
    }
  }
  onComment() {
    emitter.emit("onComment", {
      replyId: this.comment.id,
      rootId: this.comment.rootId,
      recvId: this.comment.userId,
    });
  }
  onTouchStart() {
    this.timeOutEvent = setTimeout(this.longPress, 500);
  }
  onTouchEnd(e: Event) {
    clearTimeout(this.timeOutEvent);
    if (this.timeOutEvent != 0) {
      this.onComment();
      e.preventDefault();
    }
    return false;
  }
  longPress() {
    this.timeOutEvent = 0;
    emitter.emit("more-show", { type: 7, refId: this.comment.id });
  }
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
