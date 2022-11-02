<template>
  <web-view :src=url @message="onMessage" @load="onLoad" @error="onError"/>
</template>

<script setup lang="ts">
import { WebViewProps } from "@tarojs/components";
import Taro from "@tarojs/taro";

const host = "http://127.0.0.1";

const router = Taro.getCurrentInstance().router;
let url = `${host}?platform=${Taro.getEnv()}&appUrl=${router?.path}`

const h5Url = router?.params.h5Url
if (h5Url && h5Url !== "") {
  console.log(decodeURIComponent(h5Url));
  url = `${host}/?platform=${Taro.getEnv()}&appUrl=${router?.path}#${decodeURIComponent(h5Url)}`
  console.log(url);
}
console.log("当前url",url);
function onMessage(e:WebViewProps.onMessageEventDetail) {

}

function onLoad(e:WebViewProps.onLoadEventDetail) {

}

function onError(e:WebViewProps.onErrorEventDetail) {

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
