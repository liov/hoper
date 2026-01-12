<template>
  <div>CacheStorage</div>
  <div>{{ status }}</div>
</template>

<script setup lang="ts">
import axios from "axios";
import { reactive, ref, onMounted } from "vue";
import Moment from "@h5/components/moment/Moment.vue";
import ActionMore from "@h5/components/action/More.vue";

const props = defineProps<{
  url: string | URL;
}>();

let status;

onMounted(() => {
  status = document.querySelector("#status")!;
  if ("serviceWorker" in navigator) {
    status.innerHTML += `<p>浏览器是否支持：支持</p>`;
    navigator.serviceWorker
      .register(props.url, {
        scope: "/js/",
      })
      .then(function (registration) {
        status.innerHTML += `<p>service worker是否注册成功：注册成功</p>`;

        let serviceWorker;
        if (registration.installing) {
          serviceWorker = registration.installing;
          status.innerHTML += `<p>当前注册状态：installing</p>`;
        } else if (registration.waiting) {
          serviceWorker = registration.waiting;
          status.innerHTML += `<p>当前注册状态：waiting</p>`;
        } else if (registration.active) {
          serviceWorker = registration.active;
          status.innerHTML += `<p>当前注册状态：active</p>`;
        }
        if (serviceWorker) {
          status.innerHTML += `<p>当前service worker状态：&emsp; ${serviceWorker.state}</p>`;
          serviceWorker.addEventListener("statechange", function (e) {
            status.innerHTML += `<p>状态变化为:${e.target.state}</p>`;
          });
        }
      })
      .catch(function (error) {
        console.log(error);
      });
  } else {
    status.innerHTML += `<p>浏览器是否支持：不支持</p>`;
  }
});
</script>

<style scoped lang="less"></style>
