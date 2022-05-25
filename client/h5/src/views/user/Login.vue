<template>
  <div class="login">
    <van-row type="flex" justify="center">
      <van-col span="18">
        <van-tabs @click="onClick">
          <van-tab title="登录"></van-tab>
          <van-tab title="注册"></van-tab>
        </van-tabs>
        <van-form @submit="onSubmit">
          <van-field
            v-if="type"
            v-model="username"
            name="name"
            label="用户名"
            placeholder="用户名"
            :rules="[{ required: true, message: '请填写用户名' }]"
          />
          <van-field
            v-if="type"
            name="gender"
            label="性别"
            :rules="[{ required: true, message: '请填选择性别' }]"
          >
            <template #input>
              <van-radio-group v-model="gender" direction="horizontal">
                <van-radio name="1">男</van-radio>
                <van-radio name="2">女</van-radio>
              </van-radio-group>
            </template>
          </van-field>
          <van-field
            v-if="type"
            v-model="phone"
            name="phone"
            label="手机"
            placeholder="手机"
            :rules="[
              {
                required: true,
                pattern: Validator.PhoneReg,
                message: '请填写手机',
              },
            ]"
          />
          <van-field
            v-if="type"
            v-model="mail"
            name="mail"
            label="邮箱"
            placeholder="邮箱"
            :rules="[
              {
                required: true,
                pattern: Validator.EmailReg,
                message: '请填写邮箱',
              },
            ]"
          />
          <van-field
            v-if="!type"
            v-model="account"
            name="input"
            label="账号"
            placeholder="邮箱或手机或hoper号"
            :rules="[{ required: true, message: '请填写邮箱或手机或hoper号' }]"
          />
          <van-field
            v-model="password"
            type="password"
            name="password"
            label="密码"
            placeholder="密码"
            :rules="[{ required: true, message: '请填写密码' }]"
          />
          <van-field
            v-if="type"
            type="password"
            name="password_confirm"
            v-model="password_confirm"
            label="密码确认"
            placeholder="密码"
            :rules="[
              {
                required: true,
                validator: (value, rule) => value === password,
                message: '请确保密码一致',
              },
            ]"
          />
          <Luosimao class="captcha" ref="luosimao"></Luosimao>
          <div style="margin: 16px">
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
import axios from "axios";
import Luosimao from "@/components/Luosimao.vue";
import Validator from "@/plugin/utils/validator";

@Options({
  components: { Luosimao },
})
export default class Login extends Vue {
  account = "";
  username = "";
  password = "";
  password_confirm = "";
  gender = "1";
  type = 0;
  mail = "";
  phone = "";
  Validator = Validator;
  async created() {
    if (this.$route.query.back) {
      if (!this.$store.state.user.auth) await this.$store.dispatch("getAuth");

      if (this.$store.state.user.auth)
        await this.$router.replace(`${this.$route.query.back}`);
    }
  }

  getFormValues(values: any): any {
    const res = {
      ...values,
      password: values.password,
      vCode: this.$refs.luosimao.value,
    };
    if (this.type) {
      res.gender = parseInt(this.gender);
      delete res.password_confirm;
    }
    /*    const emailReg = /^([a-zA-Z0-9]+[_.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_.]?)*[a-zA-Z0-9]+.[a-zA-Z]{2,3}$/
    const phoneReg = new RegExp('^1[0-9]{10}$')
    if (!emailReg.test(values.input) && !phoneReg.test(values.input)) {
      this.$toast.fail("请输入正确的邮箱或手机")
    }*/
    return res;
  }


  async signup(values: any) {
    const res = await axios.post("/api/v1/user", this.getFormValues(values));
    if (res.data.code == 0) {
      this.$toast.success("请前往邮箱查收激活邮件");
    }
  }

  async onSubmit(values: any) {
    if (this.type == 0)
      await this.$store.dispatch("login", this.getFormValues(values));
    else await this.$store.dispatch("signup", this.getFormValues(values)); //await this.signup(values);
    const LUOCAPTCHA = (window as any).LUOCAPTCHA;
    LUOCAPTCHA && LUOCAPTCHA.reset();
  }

  onClick(name, title) {
    this.type = name;
  }
}
</script>

<style scoped>
.login {
  margin-top: 3rem;
}

.captcha {
  display: flex;
  justify-content: center;
}
</style>
