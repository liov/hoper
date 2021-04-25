<template>
  <div class="comment">
    <span v-if="focus && comment">To:{{ user(comment.recvId).name }}</span>
    <van-field
      v-model="message"
      rows="1"
      autosize
      type="textarea"
      placeholder="请输入评论"
      :rules="[{ required: true, message: '输入内容为空' }]"
      @focus="onFocus"
      @blur="focus = false"
      ref="commentRef"
    >
      <template #button>
        <div class="button">
          <van-uploader
            v-if="focus"
            v-model="uploader"
            :max-count="1"
            :after-read="afterRead"
          />
          <van-loading v-if="loading" size="24px">上传中...</van-loading>
          <van-button size="small" type="primary" @click="onComment"
            >发送</van-button
          >
        </div>
      </template>
    </van-field>
  </div>
</template>

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
import Action from "@/components/action/Action.vue";
import axios from "axios";
import { upload } from "@/plugin/utils/upload";
import emitter from "@/plugin/emitter";
import dateTool from "@/plugin/utils/date";
import store from "@/store/index";
class Props {
  comment = prop<any>({ default: {} });
}
@Options({ components: { Action } })
export default class AddComment extends Vue.with(Props) {
  message = "";
  loading = false;
  uploader = [];
  focus = false;
  created() {
    emitter.on("onComment", (param) => {
      if (param) {
        this.comment.replyId = param.replyId;
        this.comment.rootId = param.rootId;
        this.comment.recvId = param.recvId;
      }
      this.$refs.commentRef.focus();
    });
  }
  unmounted() {
    emitter.all.delete("onComment");
  }

  async onComment() {
    if (this.message.trimStart().trimEnd().length === 0) {
      this.$toast.fail("内容为空");
      return;
    }
    const comment = {
      type: this.comment.type,
      refId: this.comment.refId,
      content: this.message,
      image: this.uploader.length > 0 ? this.uploader[0].url : "",
      replyId: this.comment.replyId,
      rootId: this.comment.rootId ? this.comment.rootId : 0,
      recvId: this.comment.recvId,
    };
    const res = await axios.post("/api/v1/action/comment", comment);
    comment.id = res.data.details.id;
    comment.userId = this.$store.state.user.auth.id;
    comment.created_at = new Date().format(dateTool.format);
    const comments = this.$store.state.content.commentCache.get(comment.rootId);
    comments.push(comment);
    this.$toast.success("评论成功");
  }
  async afterRead(file: any) {
    this.loading = true;
    file.url = await upload(file.file);
    this.loading = false;
  }
  user(id: number) {
    return this.$store.getters.getUser(id);
  }
  async onFocus() {
    if (!store.state.user.auth) {
      await this.$store.dispatch("getAuth");
    }
    if (store.state.user.auth) this.focus = true;
    else
      await this.$router.push({
        name: "Login",
        query: { back: this.$route.path },
      });
  }
}
</script>

<style scoped lang="less">
.comment {
  position: fixed;
  bottom: 47px;
  width: 100%;
}
.button {
  display: grid;
}
</style>
