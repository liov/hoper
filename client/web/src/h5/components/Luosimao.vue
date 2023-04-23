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

let value = "";
let LUOCAPTCHA = window.LUOCAPTCHA;

defineExpose({
  value,
  getValue,
});

if (!LUOCAPTCHA) {
  dynamicLoadJs("//captcha.luosimao.com/static/dist/captcha.js", () => {
    LUOCAPTCHA = window.LUOCAPTCHA;
    !window.luosimaorendered && LUOCAPTCHA.render();
    window.luosimaorendered = true;
    window.getResponse = (resp) => {
      value = resp; // resp 即验证成功后获取的值
      console.log(value);
    };
  });
}

onMounted(() => {
  LUOCAPTCHA && !window.luosimaorendered && LUOCAPTCHA.render();
  window.luosimaorendered = true;
});

function reset() {
  LUOCAPTCHA && LUOCAPTCHA.reset();
}

function getValue() {
  return value;
}
</script>

<style scoped lang="less"></style>
