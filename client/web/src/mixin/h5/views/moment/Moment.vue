<template>
  <div class="moment">
    <van-tabs
      v-model:active="active"
      @change="onChange"
      swipeable
      title-active-color="#1989fa"
      animated
      line-height="0"
    >
      <template #nav-left>
        <div style="width: 30%"></div>
      </template>
      <template #nav-right>
        <div style="width: 30%">
          <router-link to="/moment/add">
            <van-icon name="add-o" color="#1989fa" size="1.5rem" />
          </router-link>
        </div>
      </template>
      <van-tab title="关注"> 内容1 </van-tab>
      <van-tab title="推荐"> 内容2 </van-tab>
      <van-tab title="最新"> <MomentList /> </van-tab>
    </van-tabs>
    <van-row type="flex">
      <van-col span="6"> span: 6 </van-col>
      <van-col span="12">span: 12 </van-col>
      <van-col span="6">span: 6</van-col>
    </van-row>
  </div>
</template>

<script setup lang="ts">
import { showToast } from "vant";
import { defineAsyncComponent, ref } from "vue";
import MomentList from "@/mixin/h5/views/moment/List.vue";
const active = ref(2);

function onChange(index: number) {
  showToast({ message: index });
}
</script>

<style scoped lang="less">
i {
  margin: 0.7rem;
}
</style>
