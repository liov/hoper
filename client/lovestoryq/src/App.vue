<template>
  <router-view />
  <van-tabbar route>
    <van-tabbar-item replace to="/" icon="notes-o"> 瞬间 </van-tabbar-item>
    <van-tabbar-item replace to="/dairy" icon="search"> 日记 </van-tabbar-item>
    <van-tabbar-item replace to="/chat" icon="chat-o"> 聊天 </van-tabbar-item>
    <van-tabbar-item replace to="/me" icon="user-circle-o"> 我的 </van-tabbar-item>
  </van-tabbar>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";

@Options({})
export default class App extends Vue {
  created() {
    if (!this.$store.state.user.auth) {
      this.$store.dispatch("getAuth");
    }
  }
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
