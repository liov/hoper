<template>
  <van-popup
    v-model:show="show"
    position="bottom"
    teleport="#app"
    :overlay="false"
  >
    <van-field
      v-model="message"
      rows="3"
      autosize
      type="textarea"
      placeholder="请输入评论"
      :rules="[{ required: true, message: '输入内容为空' }]"
    />
    <div class="button">
      <span class="button" @click="show = false">取消</span>
      <span class="button" style="color: red" @click="onComment">确认</span>
    </div>
  </van-popup>
</template>

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
import Action from "@/components/action/Action.vue";
import axios from "axios";
class Props {
  comment = prop<any>({ default: {} });
}
@Options({ components: { Action } })
export default class AddComment extends Vue.with(Props) {
  show = false;
  message = "";
  setComment(comment) {
    this.comment.type = comment.type;
    this.comment.refId = comment.refId;
    this.comment.replyId = comment.replyId;
    this.comment.rootId = comment.rootId;
    this.comment.recvId = comment.recvId;
  }
  async onComment() {
    await axios.post("/api/v1/action/comment", {
      type: this.comment.type,
      refId: this.comment.refId,
      content: this.message,
      replyId: this.comment.replyId,
      rootId: this.comment.rootId,
      recvId: this.comment.recvId,
    });
    this.$toast.success("评论成功");
    this.show = false;
  }
}
</script>

<style scoped lang="less">
.button {
  margin: 1rem 2rem;
}
</style>
