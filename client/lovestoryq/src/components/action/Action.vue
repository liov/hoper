<template>
  <van-row>
    <van-col span="6" class="action" @click="moreShow"
      ><van-icon name="more-o"
    /></van-col>
    <van-col span="6" class="action"
      ><van-icon
        :name="content.collect ? 'star' : 'star-o'"
        :color="content.collect ? '#F6DF02' : ''"
    /></van-col>
    <van-col span="6" class="action"><van-icon name="comment-o" /></van-col>
    <van-col span="6" class="action"
      ><van-icon
        :name="content.likeId > 0 ? 'like' : 'like-o'"
        :color="content.likeId > 0 ? '#D91E46' : ''"
        @click="like(index)"
    /></van-col>
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
    emitter.emit("fav-show", { type: this.type, refId: this.content.id });
  }
  async like() {
    const api = `/api/v1/action/like`;
    const id = this.content.id;
    const likeId = this.content.likeId;
    if (likeId > 0) {
      await axios.delete(api, { data: { id: likeId } });
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

<style scoped>
.action {
  text-align: center;
}
</style>
