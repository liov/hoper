<template>
  <div class="login">
    <van-row type="flex" justify="center">
      <van-col span="18">
        <van-form @submit="onSubmit">
          <van-field
            v-model="username"
            name="username"
            label="用户名"
            placeholder="用户名"
            :rules="[{ required: true, message: '请填写用户名' }]"
          />
          <van-field
            v-model="password"
            type="password"
            name="password"
            label="密码"
            placeholder="密码"
            :rules="[{ required: true, message: '请填写密码' }]"
          />
          <div class="captcha">
            <div
                class="l-captcha"
                data-site-key="ff3498d2c6ffa1178cbf4fb6b445a8b3"
                data-width="200"
            />
          </div>
          <div style="margin: 16px;">
            <van-button round block type="primary" native-type="submit">
              提交
            </van-button>
          </div>
        </van-form>
      </van-col>
    </van-row>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import axios from "axios"

@Options({
  components: {}
})
export default class Login extends Vue {
  username = "";
  password = "";

  mounted() {
    //if (this.$route.query.email !== null) {}

    const c = document.createElement("script");
    c.type = "text/javascript";
    c.async = true;
    c.src =
      ("https:" == document.location.protocol ? "https://" : "http://") +
      "captcha.luosimao.com/static/dist/captcha.js?v=201812141420";
    const s = document.getElementsByTagName("script")[0];
    s.parentNode?.insertBefore(c, s);
    s.parentNode?.removeChild(c);
  }

  onSubmit(values: object) {
    const captcha = (document.getElementById("lc-captcha-response") as HTMLInputElement).value
    axios.post('/api/user/login',{...values,captcha:captcha})
  }
}
</script>

<style scoped>
.login {
  margin-top: 13rem;
}
.captcha{
  display:flex;
  justify-content: center;
}
</style>
