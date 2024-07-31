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
      @click-input="onFocus"
      @blur="onBlur"
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

<script setup lang="ts">
import axios from "axios";
import { upload } from "@h5/plugin/utils/upload";
import emitter from "@h5/plugin/emitter";
import dateTool from "@h5/plugin/utils/date";
import { ref, onMounted, onUnmounted, reactive, toRefs, type Ref } from "vue";
import { Toast } from "vant";
import { useRoute, useRouter } from "vue-router";
import { useUserStore } from "@h5/store/user";
import { useContentStore } from "@h5/store/content";

const props = defineProps<{
  comment: any;
}>();

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const store = useContentStore();
const message = ref("");
const loading = ref(false);
const uploader: Ref<any[]> = ref([]);
const focus = ref(false);
const commentRef = ref();
const comment = reactive(props.comment);

onMounted(() => {
  emitter.on("onComment", (param: any) => {
    if (param) {
      comment.replyId = param.replyId;
      comment.rootId = param.rootId;
      comment.recvId = param.recvId;
    }
    commentRef.value.focus();
    onFocus();
  });
});
onUnmounted(() => {
  emitter.all.delete("onComment");
});

async function onComment() {
  if (message.value.trimStart().trimEnd().length === 0) {
    Toast.fail("内容为空");
    return;
  }
  const comment: any = {
    type: props.comment.type,
    refId: props.comment.refId,
    content: message.value,
    image: uploader.value.length > 0 ? uploader.value[0].url : "",
    replyId: props.comment.replyId,
    rootId: props.comment.rootId ? props.comment.rootId : 0,
    recvId: props.comment.recvId,
  };
  const res = await axios.post("/api/v1/action/comment", comment);
  comment.id = res.data.data.id;
  comment.userId = userStore.auth.id;
  const comments = store.commentCache.get(comment.rootId);
  comments?.push(comment);
  console.log(store.commentCache);
  Toast.success("评论成功");
  message.value = "";
  focus.value = false;
}
async function afterRead(file: any) {
  loading.value = true;
  file.url = await upload(file.file);
  loading.value = false;
}
function user(id: number) {
  return userStore.getUser(id);
}
async function onFocus() {
  if (!userStore.auth) {
    await userStore.getAuth();
  }
  if (userStore.auth) focus.value = true;
  else
    await router.push({
      name: "Login",
      query: { back: route.path },
    });
}
function onBlur(e: FocusEvent) {
  if (!e.relatedTarget) focus.value = false;
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
