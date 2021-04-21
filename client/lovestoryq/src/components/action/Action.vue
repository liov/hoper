<template>
  <van-row>
    <van-col span="6" class="action" @click="moreShow"
      ><van-icon name="more-o"
    /></van-col>
    <van-col span="6" class="action"
      ><van-icon
        :name="content.collects ? 'star' : 'star-o'"
        :color="content.collects ? '#F6DF02' : ''"
        @click="favShow"
      /><span class="count">{{ content.ext.collectCount }}</span></van-col
    >
    <van-col span="6" class="action"
      ><van-icon name="comment-o" @click="commentShow" /><span class="count">{{
        content.ext.commentCount
      }}</span></van-col
    >
    <van-col span="6" class="action"
      ><van-icon
        :name="content.likeId > 0 ? 'like' : 'like-o'"
        :color="content.likeId > 0 ? '#D91E46' : ''"
        @click="like"
      /><span class="count">{{ content.ext.likeCount }}</span></van-col
    >
  </van-row>
</template>

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
import axios from "axios";
import emitter from "@/plugin/emitter";
class Props {
  content = prop<any>({ default: {} });
  readonly type = prop<number>({});
}
@Options({})
export default class Action extends Vue.with(Props) {
  moreShow() {
    emitter.emit("more-show", { type: this.type, refId: this.content.id });
  }
  favShow() {
    console.log(this.content.collect);
    emitter.emit("fav-show", {
      type: this.type,
      refId: this.content.id,
      collects: this.content.collects,
    });
  }
  commentShow() {
    emitter.emit("comment-show", {
      type: this.type,
      refId: this.type == 7 ? this.content.refId : this.content.id,
      replyId: this.type == 7 ? this.content.id : 0,
      rootId: this.type == 7 ? this.content.rootId : 0,
      recvId: this.content.userId,
    });
  }
  async like() {
    const api = `/api/v1/action/like`;
    const id = this.content.id;
    const likeId = this.content.likeId;
    if (likeId > 0) {
      await axios.delete(api + "/" + this.content.likeId);
      this.content.likeId = 0;
    } else {
      const res = await axios.post(api, {
        refId: id,
        type: this.type,
        action: 2,
      });
      this.content.likeId = res.data.details.id;
    }
  }
}
</script>

<style scoped lang="less">
.action {
  text-align: center;
  .count {
    margin: 0 2px;
  }
}
</style>
