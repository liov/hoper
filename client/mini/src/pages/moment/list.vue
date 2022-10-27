<template>

 <nut-cell round-radius="0" class="moment-list">
    <nut-list
      :height="200"
      :listData="list"
      :container-height="700"
      @scroll-bottom="onLoad"
    >
      <template v-slot="{ item }">
        <nut-cell>
          <nut-skeleton width="250px" height="200px"  title animated avatar avatarSize="50px" row="3" :loading="listConfig.loading">
            <Moment
              :moment="item"
              :user="user(item.userId)"
              :maxHeight="200"
            ></Moment>
          </nut-skeleton>
        </nut-cell>
      </template>
    </nut-list>
  </nut-cell>

</template>

<script setup lang="ts">
import axios from "@/plugins/axios";
import { onMounted, reactive, ref } from "vue";
import Moment from "@/components/moment/Moment.vue";
import { useUserStore } from "@/stores/user";
import { useContentStore } from "@/stores/content";
import type { Ref, UnwrapRef } from "vue";

const store = useContentStore();
const userStore = useUserStore();

const listConfig = reactive({
  pageNo: 1,
  pageSize: 10,
  loading: true,
  finished: false,
});

const list: Ref<UnwrapRef<any[]>> = ref(
  Array.from(new Array(6), () => {
    return {};
  })
);

onMounted(()=>onLoad())
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
    `/api/v1/moment?pageNo=${listConfig.pageNo}&pageSize=${listConfig.pageSize}`
  );
  listConfig.loading = false;
  const data = res.data.details;
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

<style lang="scss">
  .moment-list{
    padding:0
  }
  .nut-cell {
    height: 100%;
    .list-item {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 100%;
      height: 50px;
      margin-bottom: 10px;
      background-color: #f4a8b6;
      border-radius: 10px;
    }
  }
</style>
