<template>
  <div @click="reset">
    <div
      id="captcha"
      class="l-captcha"
      data-site-key="ff3498d2c6ffa1178cbf4fb6b445a8b3"
      data-width="100%"
      data-callback="getResponse"
    />
  </div>
</template>

<script setup lang="ts">
import { dynamicLoadJs } from "@/plugin/utils/script";
import { ref, onMounted } from "vue";

let second = false;
const value = ref("");

if (!window.LUOCAPTCHA) {
  dynamicLoadJs("//captcha.luosimao.com/static/dist/captcha.js");
} else second = true;
window.LUOCAPTCHA.getResponse = (resp) => {
  value.value = resp; // resp 即验证成功后获取的值
};

onMounted(() => {
  if (second) window.LUOCAPTCHA.render();
});

function reset() {
  const LUOCAPTCHA = window.LUOCAPTCHA;
  LUOCAPTCHA && LUOCAPTCHA.reset();
}
</script>

<style scoped lang="less"></style>
