<template>
  <div>
    <div>Web Woker</div>
    <div id="result" />
    <a-button @click="stopWorker"> 停止 </a-button>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from "vue";

let w;
onMounted(() => {
  startWorker();
  const vm = this;
  setTimeout(function () {
    stopWorker();
  }, 3000);
});

function startWorker() {
  if (typeof Worker !== "undefined") {
    if (typeof w === "undefined") {
      w = new Worker("../js/demo_workers.js");
    }
    w.onmessage = function (event) {
      document.getElementById("result")!.innerHTML = event.data;
    };
  } else {
    document.getElementById("result")!.innerHTML =
      "抱歉，你的浏览器不支持 Web Workers...";
  }
}
function stopWorker() {
  w.terminate();
  w = undefined;
}
</script>

<style scoped></style>
