<template>
  <RouterView v-slot="{ Component }">
    <template v-if="Component">
      <Transition name="fade">
        <KeepAlive>
          <Suspense>
            <!-- main content -->
            <component :is="Component"></component>

            <!-- loading state -->
            <template #fallback> Loading... </template>
          </Suspense>
        </KeepAlive>
      </Transition>
    </template>
  </RouterView>
  <van-tabbar route>
    <van-tabbar-item replace to="/" icon="notes-o"> 瞬间 </van-tabbar-item>
    <van-tabbar-item replace to="/dairy" icon="search"> 日记 </van-tabbar-item>
    <van-tabbar-item replace to="/chat" icon="chat-o"> 聊天 </van-tabbar-item>
    <van-tabbar-item replace to="/me" icon="user-circle-o"
      >我的</van-tabbar-item
    >
  </van-tabbar>
</template>

<script setup lang="ts">
import { RouterView, useRouter } from "vue-router";
import { useGlobalStore } from "@/store/global";
import { Platform } from "@/model/const";
import { dynamicLoadJs } from "@/plugin/utils/script";
import wxenv from "@/plugin/platform/weixin";

const router = useRouter();
const store = useGlobalStore();
switch (router.currentRoute.value.query.platform) {
  case Platform.Weapp:
    store.platform = Platform.Weapp;
    if (!wxenv.wx) {
      dynamicLoadJs("//res.wx.qq.com/open/js/jweixin-1.3.2.js", () => {
        wxenv.wx = window.wx;
      });
    }
}

if (wxenv.IsWeappPlatform() && !wxenv.wx) {
  dynamicLoadJs("//res.wx.qq.com/open/js/jweixin-1.3.2.js", () => {
    wxenv.wx = window.wx;
  });
}
</script>

<style lang="less">
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

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
