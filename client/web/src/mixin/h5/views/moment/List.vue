<template>
  <van-pull-refresh
    v-model="pullDown.refreshing"
    :success-text="pullDown.successText"
    @refresh="onRefresh"
  >
    <van-list
      v-model:loading="listConfig.loading"
      :finished="listConfig.finished"
      finished-text="没有更多了"
      @load="onLoad"
    >
      <van-cell v-for="item in list" :key="item.id">
        <van-skeleton title avatar round :row="3" :loading="listConfig.loading">
          <Moment
            v-if="show"
            :moment="item"
            :user="user(item.userId)"
            :maxHeight="200"
          ></Moment>
        </van-skeleton>
      </van-cell>
    </van-list>
  </van-pull-refresh>
  <ActionMore v-if="userStore.auth" key="moment-list"></ActionMore>
</template>

<script setup lang="ts">
import axios from "axios";
import { onMounted, reactive, ref } from "vue";
import Moment from "@/mixin/h5/components/moment/Moment.vue";
import ActionMore from "@/mixin/h5/components/action/More.vue";
import { useUserStore } from "@/mixin/store/user";
import { useContentStore } from "@/mixin/store/content";
import type { Ref } from "vue";
import { momentList } from "@/mixin/service/moment";

const store = useContentStore();
const userStore = useUserStore();

const listConfig = reactive({
  pageNo: 1,
  pageSize: 10,
  loading: false,
  finished: false,
});

const list: Ref<any[]> = ref(
  Array.from(new Array(listConfig.pageSize), () => {
    return {};
  }),
);

const pullDown = reactive({
  refreshing: false,
  successText: "刷新成功",
});
const show = ref(false);

function user(id: number) {
  return userStore.getUser(id);
}

async function onLoad() {
  listConfig.finished = false;
  // 异步更新数据
  const res = await axios.get(
    `/api/v1/moment?pageNo=${listConfig.pageNo}&pageSize=${listConfig.pageSize}`,
  );

  listConfig.loading = false;
  const data = res.data.data;
  if (!data || !data.list) {
    listConfig.finished = true;
    return;
  }
  if (listConfig.pageNo == 1) {
    list.value = data.list;
  } else {
    list.value = list.value.concat(data.list);
  }
  console.log(list);
  userStore.appendUsers(data.users);
  show.value = true;
  listConfig.pageNo++;
  if (data.list.length < listConfig.pageSize) listConfig.finished = true;
}

function onRefresh() {
  pullDown.refreshing = true;
  listConfig.pageNo = 1;
  onLoad().catch(() => {
    pullDown.successText = "刷新失败";
  });
  pullDown.refreshing = false;
}
</script>

<style scoped lang="less"></style>
