<template>
  <van-config-provider :theme="theme">
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
      <van-tabbar-item replace to="/diary" icon="search">
        日记
      </van-tabbar-item>
      <van-tabbar-item replace to="/chat" icon="chat-o"> 聊天 </van-tabbar-item>
      <van-tabbar-item replace to="/me" icon="user-circle-o"
        >我的</van-tabbar-item
      >
    </van-tabbar>
  </van-config-provider>
</template>

<script setup lang="ts">
import { RouterView, useRouter } from "vue-router";
import { useGlobalStore } from "@/mixin/store/global";
import { Platform } from "@/mixin/model/const";
import wxenv from "@/utils/platform/weixin";
import { parseQueryString } from "@hopeio/utils/browser";
import "@types/grpc_custom_status";
import { ref } from "vue";
import type { Ref } from "vue";
import type { ConfigProviderTheme } from "vant";
//  Type 'string' is not assignable to type 'ConfigProviderTheme | undefined'
const light: ConfigProviderTheme = "light";
const theme = ref(light);

const store = useGlobalStore();
console.log("url:", window.location.href);
const queryParams = parseQueryString();
console.log(queryParams);
if (queryParams.platform) {
  switch (queryParams.platform.toUpperCase()) {
    case Platform.Weapp:
      store.platform = Platform.Weapp;
      if (!window.wx) {
        wxenv.loadwxSDK();
        console.log(window.wx);
      }
      break;
    case Platform.App:
      store.platform = Platform.App;
      break;
  }
} else {
  store.platform = Platform.H5;
}

if (wxenv.IsWeappPlatform() && !window.wx) {
  wxenv.loadwxSDK();
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
