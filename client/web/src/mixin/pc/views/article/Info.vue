<template>
  <a-row>
    <a-col :span="3" />
    <a-col :span="18">
      <div
        style="
          height: 60px;
          line-height: 60px;
          margin: 10px;
          text-align: center;
          font-size: 26px;
          color: rgba(0, 0, 0, 0.65);
        "
      >
        {{ article.title }}
      </div>
      <a-divider />
      <div v-html="article.html_content" />
      <a-divider />
      <div style="margin: 0 auto">
        <span @click="star">
          <star-outlined style="margin-right: 8px" />
          {{ article.collect_count }}
        </span>
        <span @click="like">
          <like-outlined style="margin-right: 8px" />
          {{ article.like_count }}
        </span>
        <span @click="comment">
          <message-outlined style="margin-right: 8px" />
          {{ article.comment_count }}
        </span>
        <a-divider type="vertical" />
        <a-tag
          v-for="(subitem, subindex) in article.categories"
          :key="subindex"
        >
          {{ subitem.name }}
        </a-tag>
        <a-divider type="vertical" />
        <a-checkable-tag
          v-for="(subitem, subindex) in article.tags"
          :key="subindex"
          :color="color[subindex]"
        >
          {{ subitem.name }}
        </a-checkable-tag>
        <a-divider type="vertical" />
        <a-rate :default-value="2.5" allow-half />
      </div>
    </a-col>
    <a-col :span="3" />
  </a-row>
</template>

<script setup lang="ts">
import { onBeforeUpdate, onMounted, onUpdated, reactive, ref } from "vue";
import type { Ref } from "vue";
import { tagColor } from "@/mixin/pc/views/article/const";
import ArticleClient from "@/mixin/service/article";
import { useRoute, useRouter } from "vue-router";
import {
  StarOutlined,
  LikeOutlined,
  MessageOutlined,
  EditOutlined,
} from "@ant-design/icons-vue";
import axios from "axios";
import { message } from "ant-design-vue";
import type MarkdownIt from "markdown-it";
let md: MarkdownIt;
const route = useRoute();
const color: Ref<string[]> = ref(tagColor);
const article: Ref<any> = ref({});
article.value = await ArticleClient.info(route.params.id);
const content = ref(article.value.html_content);

if (article.value.content_type == 0) {
  const MarkDownIt = await import("markdown-it");
  md = MarkDownIt.default();
  article.value.html_content = md.render(article.value.content);
}

onBeforeUpdate(async () => {
  article.value = await ArticleClient.info(route.params.id);
  if (!md && article.value.content_type == 0) {
    const MarkDownIt = await import("markdown-it");
    md = MarkDownIt.default();
  }
  article.value.html_content = md.render(article.value.content);
});

async function star() {}
function like() {}
function comment() {}
</script>

<style scoped lang="less"></style>
