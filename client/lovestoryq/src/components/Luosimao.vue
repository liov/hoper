<template>
  <div @click="reset">
    <div
      class="l-captcha"
      data-site-key="ff3498d2c6ffa1178cbf4fb6b445a8b3"
      data-width="100%"
      data-callback="getResponse"
    />
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";

@Options({
  components: {}
})
export default class Luosimao extends Vue {
  created() {
    const c = document.createElement("script");
    c.type = "text/javascript";
    c.async = true;
    c.src =
      ("https:" == document.location.protocol ? "https://" : "http://") +
      "captcha.luosimao.com/static/dist/captcha.js?v=201812141420";
    const s = document.getElementsByTagName("script")[0];
    s.parentNode?.insertBefore(c, s);
    //s.parentNode?.removeChild(c);
  }

  mounted() {
    (window as any).getResponse = resp => {
      if (this.captcha() === resp) this.$emit("success", resp);
      else this.reset();
    };
  }

  captcha = () =>
    (document.getElementById("lc-captcha-response") as HTMLInputElement).value;

  reset() {
    const LUOCAPTCHA = (window as any).LUOCAPTCHA;
    LUOCAPTCHA && LUOCAPTCHA.reset();
    console.log("执行了");
  }
}
</script>

<style scoped lang="less"></style>
