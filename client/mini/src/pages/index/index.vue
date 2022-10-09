<template>
  <div>
  <Moment v-if="activeName==='moment'"></Moment>
  <nut-tabbar @tab-switch="tabSwitch" bottom safeAreaInsetBottom v-model:visible="activeName">
    <nut-tabbar-item icon="image" tab-title="瞬间" name="moment">
    </nut-tabbar-item>
    <nut-tabbar-item  icon="search" tab-title="日记" name="search"></nut-tabbar-item>
    <nut-tabbar-item  icon="message" tab-title="聊天" name="message"></nut-tabbar-item>
    <nut-tabbar-item  icon="my" tab-title="我的" name="my"></nut-tabbar-item>
  </nut-tabbar>
  </div>
</template>

<script setup lang="ts">

import { useUserStore } from "@/stores/user";
import articleIcon  from  '@/assets/icon/article.svg'
import Moment from '@/pages/moment/moment.vue'
import { ref } from "vue";

const activeName = ref("moment");
function tabSwitch(item:any, index:number) {
  console.log(item, index);
}

</script>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;
  }
}
</style>
