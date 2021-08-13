<template>
  <van-row>
    <van-col span="6" class="action" @click="moreShow"
      ><van-icon name="more-o"
    /></van-col>
    <van-col span="6" class="action"
      ><van-icon
        :name="
          content.collects && content.collects.length > 0 ? 'star' : 'star-o'
        "
        :color="
          content.collects && content.collects.length > 0 ? '#F6DF02' : ''
        "
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
import { jump } from "@/router/utils";

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
    if (!this.content.collects) {
      this.content.collects = [];
    }
    emitter.emit("fav-show", {
      type: this.type,
      refId: this.content.id,
      collects: this.content.collects,
    });
    console.log(this.content.collects);
  }
  commentShow() {
    jump(this.$route.path, this.type, this.content);
  }
  async like() {
    const api = `/api/v1/action/like`;
    const id = this.content.id;
    const likeId = this.content.likeId;
    if (likeId > 0) {
      await axios.delete(api + "/" + this.content.likeId);
      this.content.likeId = 0;
    } else {
      const { data } = await axios.post(api, {
        refId: id,
        type: this.type,
        action: 2,
      });
      this.content.likeId = data.details.id;
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
