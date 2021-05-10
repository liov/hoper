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

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import { dynamicLoadJs } from "@/plugin/utils/script";
@Options({
  components: {},
})
export default class Luosimao extends Vue {
  second = false;
  value = "";
  created() {
    if (!this.$store.state.LUOCAPTCHA) {
      dynamicLoadJs("//captcha.luosimao.com/static/dist/captcha.js");
      this.$store.commit("setCaptcha", (window as any).LUOCAPTCHA);
    } else this.second = true;
    window.getResponse = (resp) => {
      this.value = resp; // resp 即验证成功后获取的值
    };
  }
  mounted() {
    if (this.second) this.$store.state.LUOCAPTCHA.render();
  }

  reset() {
    const LUOCAPTCHA = (window as any).LUOCAPTCHA;
    LUOCAPTCHA && LUOCAPTCHA.reset();
  }
}
</script>

<style scoped lang="less"></style>
