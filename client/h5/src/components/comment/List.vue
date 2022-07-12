<template>
  <div>
    <van-list
      v-if="list.length > 0"
      v-model:loading="loading"
      :finished="finished"
      finished-text="没有更多了"
      @load="onLoad"
    >
      <van-cell v-for="item in list" :key="item.id">
        <template #default>
          <van-skeleton title avatar round :row="3" :loading="loading">
            <Comment :comment="item" :user="user(item.userId)"></Comment>
          </van-skeleton>
        </template>
      </van-cell>
    </van-list>
    <ActionMore key="comment-list"></ActionMore>
  </div>
</template>

<script setup lang="ts">
import axios from "axios";
import { reactive, type Ref, ref, type UnwrapNestedRefs } from "vue";
import Comment from "@/components/comment/Comment.vue";
import ActionMore from "@/components/action/More.vue";
import { useUserStore } from "@/store/user";
import { useContentStore } from "@/store/content";
import type { UnwrapRef } from "vue";

const props = withDefaults(
  defineProps<{
    type: number;
    refId: string;
    rootId?: number;
  }>(),
  { rootId: 0 }
);

const userStore = useUserStore();
const store = useContentStore();

const active = 0;
const loading = ref(false);
const finished = ref(false);
const pageNo = ref(1);
const pageSize = ref(10);
const list: Ref<UnwrapRef<any[]>> = ref([]);

const pullDown = reactive({
  refreshing: false,
  successText: "刷新成功",
});

onLoad();

//mounted() {}
function user(id: number) {
  return userStore.getUser(id);
}

async function onLoad() {
  finished.value = false;
  // 异步更新数据
  const res = await axios.get(
    `/api/v1/action/comment?type=${props.type}&refId=${props.refId}&rootId=${props.rootId}&pageNo=${pageNo.value}&pageSize=${pageSize.value}`
  );
  loading.value = false;
  const data = res.data.details;
  if (!data || !data.list) {
    store.commentCache.set(props.rootId, list);
    finished.value = true;
    return;
  }
  if (pageNo.value == 1) {
    list.value = data.list;
  } else {
    list.value = list.value.concat(data.list);
  }
  store.commentCache.set(props.rootId, list);

  userStore.appendUsers(data.users);
  pageNo.value++;
  if (data.list.length < pageSize.value) finished.value = true;
}
</script>

<style scoped lang="less"></style>
