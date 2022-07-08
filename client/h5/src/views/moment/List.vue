<template>
  <div>
    <van-pull-refresh
      v-model="pullDown.refreshing"
      :success-text="pullDown.successText"
      @refresh="onRefresh"
    >
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-cell v-for="item in list" :key="item.id">
          <template #default>
            <van-skeleton title avatar round :row="3" :loading="loading">
              <Moment
                v-if="show"
                :moment="item"
                :user="user(item.userId)"
                :maxHeight="200"
              ></Moment>
            </van-skeleton>
          </template>
        </van-cell>
      </van-list>
    </van-pull-refresh>
    <ActionMore v-if="userStore.auth" key="moment-list"></ActionMore>
  </div>
</template>

<script setup lang="ts">
import axios from "axios";
import { reactive, ref } from "vue";
import Moment from "@/components/moment/Moment.vue";
import ActionMore from "@/components/action/More.vue";
import { useUserStore } from "@/store/user";
import { useContentStore } from "@/store/content";
import type { UnwrapNestedRefs } from "vue";

const store = useContentStore();
const userStore = useUserStore();
const active = 0;
const loading = ref(false);
const finished = ref(false);
const pageNo = ref(1);
const pageSize = ref(10);
const list: UnwrapNestedRefs<any> = ref(
  Array.from(new Array(pageSize), () => {
    return {};
  })
);

const pullDown = reactive({
  refreshing: false,
  successText: "刷新成功",
});
const show = ref(false);

//mounted() {}
function user(id: number) {
  return userStore.getUser(id);
}

async function onLoad() {
  finished.value = false;
  // 异步更新数据
  const res = await axios.get(
    `/api/v1/moment?pageNo=${pageNo.value}&pageSize=${pageSize.value}`
  );
  loading.value = false;
  const data = res.data.details;
  if (!data || !data.list) {
    finished.value = true;
    return;
  }
  if (pageNo.value == 1) {
    list.value = data.list;
  } else {
    list.value = list.value.concat(data.list);
  }
  console.log(list);
  userStore.appendUsers(data.users);
  show.value = true;
  pageNo.value++;
  if (data.list.length < pageSize.value) finished.value = true;
}

function onRefresh() {
  pullDown.refreshing = true;
  pageNo.value = 1;
  onLoad().catch(() => {
    pullDown.successText = "刷新失败";
  });
  pullDown.refreshing = false;
}
</script>

<style scoped lang="less"></style>
