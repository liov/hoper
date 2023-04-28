<template>
  <a-row :gutter="24">
    <a-col :span="12">
      <router-link to="/article/add"> ss </router-link>
    </a-col>
    <a-col v-if="!isLogin" :span="12">
      <a-form
        :model="user"
        name="login"
        :label-col="formItemLayout.labelCol"
        :wrapper-col="formItemLayout.wrapperCol"
        autocomplete="off"
        @finish="onFinish"
      >
        <a-form-item
          label=""
          :label-col="{ span: 5, offset: 5 }"
          :wrapper-col="{ span: 6, offset: 5 }"
        >
          <a-radio-group v-model:value="formType" @change="handleChange">
            <a-radio-button value="login"> 登录 </a-radio-button>
            <a-radio-button value="signup"> 注册 </a-radio-button>
          </a-radio-group>
        </a-form-item>

        <a-form-item
          v-show="formType === 'signup'"
          label="用户名"
          :rules="[
            { required: formType === 'signup', message: '请输入用户名!' },
          ]"
        >
          <a-input v-model:value="user.name" placeholder="请输入用户名！" />
        </a-form-item>

        <a-form-item
          v-show="formType === 'login'"
          label="邮箱/手机"
          :rules="[
            { required: formType === 'login', message: '请输入邮箱或手机!' },
          ]"
        >
          <a-input
            v-model:value="user.input"
            type="email"
            placeholder="输入邮箱或手机号！"
          />
        </a-form-item>
        <a-form-item
          v-show="formType === 'signup'"
          :rules="[{ required: formType === 'signup', message: '请输入邮箱!' }]"
          label="邮箱"
        >
          <a-input
            v-model:value="user.email"
            type="email"
            placeholder="请输入邮箱！"
          />
        </a-form-item>
        <a-form-item
          label="密码"
          :rules="[{ required: true, message: '请输入密码!' }]"
        >
          <a-input
            v-model:value="user.password"
            type="password"
            placeholder="请输入密码！"
          />
        </a-form-item>

        <a-form-item
          v-show="formType === 'signup'"
          label="性别"
          required
          :rules="[{ required: formType === 'signup', message: '请选择性别!' }]"
        >
          <a-radio-group v-model:value="user.sex">
            <a-radio-button value="男"> 男 </a-radio-button>
            <a-radio-button value="女"> 女 </a-radio-button>
          </a-radio-group>
        </a-form-item>

        <a-form-item
          v-show="formType === 'signup'"
          :rules="[
            { required: formType === 'signup', message: '请输入手机号!' },
          ]"
          label="手机号"
        >
          <a-input
            v-model:value="user.phone"
            type="phone"
            placeholder="请输入手机号!"
          />
        </a-form-item>
        <a-form-item
          :label-col="{ span: 5, offset: 5 }"
          :wrapper-col="{ span: 6, offset: 5 }"
        >
          <div
            class="l-captcha"
            data-site-key="ff3498d2c6ffa1178cbf4fb6b445a8b3"
            data-width="200"
          />
        </a-form-item>
        <a-form-item
          :label-col="{ span: 6 }"
          :wrapper-col="{ span: 8, offset: 5 }"
        >
          <a-button type="primary" html-type="submit">
            {{ formType === "login" ? "登录" : "注册" }}
          </a-button>
        </a-form-item>
      </a-form>
    </a-col>
    <a-col v-if="isLogin" :span="12">
      <router-link to="/user/self">
        <a-button> 个人信息 </a-button>
      </router-link>
      <a-button @click="logout"> 注销 </a-button>
    </a-col>
  </a-row>
</template>

<script setup lang="ts">
import { useUserStore } from "@pc/store/user";
import { reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import axios from "axios";
import { message, Form } from "ant-design-vue";

const formItemLayout = {
  labelCol: { span: 5 },
  wrapperCol: { span: 7 },
};
const formTailLayout = {
  labelCol: { span: 5 },
  wrapperCol: { span: 8, offset: 6 },
};
const formType = ref("login");
const isLogin = ref(false);

const useForm = Form.useForm;
const user = reactive({
  name: "",
  email: "",
  input: "",
  password: "",
  phone: "",
  sex: "男",
});
const route = useRoute();
const router = useRouter();

const userStore = useUserStore();
if (userStore.auth !== null) {
  isLogin.value = true;
}

if (route.query.email !== null) {
  user.email = route.query.email as string;
}

const c = document.createElement("script");
c.type = "text/javascript";
c.async = true;
c.src =
  ("https:" == document.location.protocol ? "https://" : "http://") +
  "captcha.luosimao.com/static/dist/captcha.js?v=201812141420";
const s = document.getElementsByTagName("script")[0];
s.parentNode?.insertBefore(c, s);
s.parentNode?.removeChild(c);
function handleChange(e) {
  formType.value = e.target.value;
  const emailReg = new RegExp(
    "^([a-zA-Z0-9]+[_.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_.]?)*[a-zA-Z0-9]+.[a-zA-Z]{2,3}$"
  );
  const phoneReg = /^1[0-9]{10}$/;
  if (emailReg.test(user.input)) {
    user.email = user.input;
  } else if (phoneReg.test(user.input)) {
    user.phone = user.input;
  }
  /* this.$nextTick(() => {
        this.user.validateFields(['password'], { force: true })
      }) */
}
// 尊重语法糖，抛弃恶心的回调
function onFinish(values: any) {
  axios
    .post(`/api/user/` + formType.value, {
      ...values,
      luosimao: document.getElementsByName("luotest_response")[0].innerText,
    })
    .then(({ data }) => {
      // success
      if (data.code === 200) {
        if (data.msg === "登录成功") {
          localStorage.setItem("token", data.token);
          userStore.auth = data.data;
          userStore.token = data.token;
          localStorage.setItem("user", data.data.id);
          message.info("登录成功");
          if (route.query.callbackUrl) {
            router.replace(route.query.callbackUrl as string);
          } else {
            router.replace("/");
          }
        } else if (data.msg === "注册成功") {
          message.info("注册成功，请到邮箱激活");
        }
        // vm.$router.replace('/')
      } else if (data.msg === "账号未激活") {
        message.warning(data.user.email);
      } else {
        message.error(data.msg);
      }
    })
    .catch(function (err) {
      message.error(err);
    });
}
async function logout() {
  const { data } = await axios.get("/api/user/logout");
  localStorage.removeItem("token");
  localStorage.removeItem("user");
  userStore.auth = null;
  userStore.token = "";
  message.info(data.msg);
  isLogin.value = false;
}
</script>

<style scoped></style>
