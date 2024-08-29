<template>
  <a-config-provider :locale="zh_CN">
    <a-layout>
      <a-layout-header
        :style="{
          padding: 0,
        }"
      >
        <div class="menu">
          <a-menu
            v-model:selectedKeys="current"
            :theme="theme"
            mode="horizontal"
            :style="{ lineHeight: '64px', justifyContent: 'center' }"
          >
            <a-menu-item key="home">
              <router-link to="/">
                <home-outlined />
                主页
              </router-link>
            </a-menu-item>
            <a-menu-item key="file-text">
              <router-link to="/article">
                <file-text-outlined />
                博客
              </router-link>
            </a-menu-item>

            <a-menu-item key="message">
              <router-link to="/chat/v2">
                <message-outlined />
                聊天
              </router-link>
            </a-menu-item>

            <a-menu-item key="picture">
              <router-link to="/moment">
                <picture-outlined />
                瞬间
              </router-link>
            </a-menu-item>
            <a-menu-item key="user">
              <router-link to="/user/login">
                <user-outlined />
                <span v-if="userStore.auth">注销</span>
                <span v-else>登录</span>
              </router-link>
            </a-menu-item>
            <!--<a-menu-item key="book">
            <nuxt-link to="/diary">
              <a-icon type="book" /> 日记
            </nuxt-link>
          </a-menu-item>-->
            <a-sub-menu>
              <template #title
                ><span class="submenu-title-wrapper"
                  ><setting-outlined />设置</span
                ></template
              >
              <a-menu-item-group title="初始化">
                <a-menu-item key="setting:1">
                  <router-link to="/api/init"> 数据库初始化 </router-link>
                </a-menu-item>
                <a-menu-item key="setting:2"> 设置初始化 </a-menu-item>
              </a-menu-item-group>
              <a-menu-item-group title="危险操作">
                <a-menu-item key="setting:3"> 重启 </a-menu-item>
                <a-menu-item key="setting:4"> 关闭 </a-menu-item>
              </a-menu-item-group>
            </a-sub-menu>
            <a-menu-item key="app" disabled>
              <appstore-outlined />
              Hoper
            </a-menu-item>
            <a-menu-item key="app" disabled>
              <a-switch
                :default-checked="false"
                @change="changeTheme"
              ></a-switch>
              主题
            </a-menu-item>
          </a-menu>
        </div>
      </a-layout-header>
      <a-layout-content>
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
      </a-layout-content>
      <a-layout-footer style="text-align: center">
        hoper ©2019 Created by JYB
        <a href="https://beian.miit.gov.cn/" target="_blank"
          >晋ICP备18012261-1号</a
        >
        <div style="width: 300px; margin: 0 auto; padding: 20px 0">
          <a
            target="_blank"
            href="http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=44030602007392"
            style="
              display: inline-block;
              text-decoration: none;
              height: 20px;
              line-height: 20px;
            "
            ><img src="@/mixin/pc/assets/beian.png" style="float: left" />
            <p
              style="
                float: left;
                height: 20px;
                line-height: 20px;
                margin: 0px 0px 0px 5px;
                color: #939393;
              "
            >
              粤公网安备 44030602007392号
            </p></a
          >
        </div>
      </a-layout-footer>
    </a-layout>
  </a-config-provider>
</template>

<script setup lang="ts">
import { RouterView, useRoute, useRouter } from "vue-router";
import { useGlobalStore } from "@/mixin/store/global";
import { useUserStore } from "@/mixin/store/user";
import { parseQueryString } from "@/utils/location";
import "@types/grpc_custom_status";
import zh_CN from "ant-design-vue/es/locale/zh_CN";
import { onDeactivated, onUnmounted, ref } from "vue";
import {
  FileTextOutlined,
  UserOutlined,
  PictureOutlined,
  SettingOutlined,
  MessageOutlined,
  HomeOutlined,
  AppstoreOutlined,
} from "@ant-design/icons-vue";

const store = useGlobalStore();
const userStore = useUserStore();
console.log("url:", window.location.href);
const queryParams = parseQueryString();
console.log(queryParams);
const route = useRoute();

const menu = ref([
  { title: "主页", link: "/", icon: "home" },
  {
    title: "文章",
    link: "/article",
    icon: "file-text",
  },
  { title: "", link: "", icon: "" },
]);
const current = ref(["main"]);
const theme = ref("light");
const style: any = ref({ backgroundColor: "#fff" });

onUnmounted(() => sessionStorage.setItem("back_url", route.path));
function changeTheme(checked: boolean) {
  theme.value = checked ? "dark" : "light";
  style.value = checked
    ? { height: "48px" }
    : { backgroundColor: "#fff", height: "48px" };
}
</script>

<style lang="less" scoped>
.menu {
  width: 100%;
  text-align: center;
  top: 0;
  z-index: 1;
}
</style>
