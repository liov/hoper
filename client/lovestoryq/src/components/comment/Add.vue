<template></template>

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
import Action from "@/components/action/Action.vue";
import axios from "axios";
class Props {
  comment = prop<any>({ default: {} });
}
@Options({ components: { Action } })
export default class AddComment extends Vue.with(Props) {
  message = "";
  async onComment() {
    await axios.post("/api/v1/action/comment", {
      type: this.comment.type,
      refId: this.comment.refId,
      content: this.message,
      replyId: this.comment.replyId,
      rootId: this.comment.rootId,
      recvId: this.comment.recvId,
    });
  }
}
</script>

<style scoped></style>
