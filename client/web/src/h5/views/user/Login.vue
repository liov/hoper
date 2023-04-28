<template>
  <div class="login">
    <van-row type="flex" justify="center">
      <van-col span="18">
        <van-tabs @click-tab="onClickTab">
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

<script setup lang="ts">
import { nextTick, onMounted, ref } from "vue";
import axios from "axios";
import Luosimao from "@h5/components/Luosimao.vue";
import Validator from "@h5/plugin/utils/validator";
import { useRoute, useRouter } from "vue-router";
import { useUserStore } from "@h5/store/user";
import { Dialog, Toast } from "vant";
import { useGlobalStore } from "@h5/store/global";
import { Platform } from "@h5/model/const";

const globalState = useGlobalStore();
onMounted(() => {
  if (globalState.platform === Platform.Weapp) {
    Dialog.confirm({
      title: "微信登录",
      message: "使用微信登录",
    })
      .then(() => {
        // on confirm
        window.wx.miniProgram.navigateTo({
          url: `/pages/user/login?h5Url=${encodeURIComponent(route.fullPath)}`,
        });
      })
      .catch(() => {
        // on cancel
      });
  }
});

const account = ref("");
const username = ref("");
const password = ref("");
const password_confirm = ref("");
const gender = ref("1");
const type = ref(0);
const mail = ref("");
const phone = ref("");
const router = useRouter();
const route = useRoute();
const store = useUserStore();
const luosimao: any = ref(null);

if (route.query.back) {
  if (!store.auth) await store.getAuth();

  if (store.auth) await router.replace(`${route.query.back}`);
}

function getFormValues(values: any): any {
  nextTick(() => {
    console.log(luosimao.value.value);
  });
  const res = {
    ...values,
    password: values.password,
    vCode: luosimao.value.getValue(),
  };
  if (type.value) {
    res.gender = parseInt(gender.value);
    delete res.password_confirm;
  }
  /*    const emailReg = /^([a-zA-Z0-9]+[_.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_.]?)*[a-zA-Z0-9]+.[a-zA-Z]{2,3}$/
    const phoneReg = new RegExp('^1[0-9]{10}$')
    if (!emailReg.test(values.input) && !phoneReg.test(values.input)) {
      this.$toast.fail("请输入正确的邮箱或手机")
    }*/
  return res;
}

async function signup(values: any) {
  const res = await axios.post("/api/v1/user", getFormValues(values));
  if (res.data.code == 0) {
    Toast.success("请前往邮箱查收激活邮件");
  }
}

async function onSubmit(values: any) {
  if (type.value == 0) await store.login(getFormValues(values));
  else await store.signup(getFormValues(values)); //await this.signup(values);
  const LUOCAPTCHA = window.LUOCAPTCHA;
  LUOCAPTCHA && LUOCAPTCHA.reset();
}

function onClickTab({ name, title }) {
  console.log(name, title);
  type.value = name;
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
