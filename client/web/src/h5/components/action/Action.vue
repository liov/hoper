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

<script setup lang="ts">
import axios from "axios";
import emitter from "@/plugin/emitter";
import { jump } from "@/router/utils";
import { reactive } from "vue";
import { useRoute } from "vue-router";

const props = defineProps<{
  content: any;
  readonly type: number;
}>();

const route = useRoute();

const content = reactive(props.content);

function moreShow() {
  emitter.emit("more-show", { type: props.type, refId: content.id });
}
function favShow() {
  if (!props.content.collects) {
    content.collects = [];
  }
  emitter.emit("fav-show", {
    type: props.type,
    refId: content.id,
    collects: content.collects,
  });
  console.log(content.collects);
}
function commentShow() {
  jump(route.path, props.type, content);
}
async function like() {
  const api = `/api/v1/action/like`;
  const id = content.id;
  const likeId = content.likeId;
  if (likeId > 0) {
    await axios.delete(api + "/" + content.likeId);
    content.likeId = 0;
  } else {
    const { data } = await axios.post(api, {
      refId: id,
      type: props.type,
      action: 2,
    });
    content.likeId = data.details.id;
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
