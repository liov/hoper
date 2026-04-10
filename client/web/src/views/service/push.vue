<template>
  <div>Push</div>
</template>

<script setup lang="ts">
import Push from "push.js";
import { onMounted } from "vue";

onMounted(() => {
  const promise = Push.create("Hello world!", {
    body: "How's it hangin'?",
    icon: "/icon.png",
    timeout: 4000,
    onClick: function () {
      window.focus();
      this.close();
    },
  });
  promise.then(function (notification) {
    notification.close();
  });
  Push.clear();
});
</script>

<style scoped></style>
