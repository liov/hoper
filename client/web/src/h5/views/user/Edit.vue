<template>
  <van-cell-group>
    <van-row>
      <van-col span="10">
        <van-uploader
          :deletable="false"
          :max-count="1"
          :preview-full-image="false"
          :after-read="afterRead"
        >
          <template #default>
            <van-image
              width="100"
              height="100"
              :src="staticDir + user.avatarUrl"
              fit="cover"
            />
          </template>
        </van-uploader>
        <van-loading v-if="loading" size="24px">上传中...</van-loading>
        <p>点击可更换头像</p>
      </van-col>
      <van-col span="14">
        <van-field
          label-width="50"
          v-model="user.name"
          label="昵称"
          :placeholder="user.name"
        />
        <van-field label-width="50" name="radio" label="性别">
          <template #input>
            <van-radio-group v-model="user.gender" direction="horizontal">
              <van-radio name="1">男</van-radio>
              <van-radio name="2">女</van-radio>
            </van-radio-group>
          </template>
        </van-field>
        <van-field
          label-width="50"
          v-model="user.birthday"
          label="生日"
          @click="show = true"
        />
      </van-col>
    </van-row>
    <van-field
      v-model="user.signature"
      label="个性签名"
      :placeholder="user.signature"
    />
    <van-field
      v-model="user.intro"
      rows="1"
      autosize
      label="个人简介"
      type="textarea"
      :placeholder="user.intro"
    />
  </van-cell-group>
  <van-popup v-if="show" v-model:show="show" position="bottom" teleport="#app">
    <van-datetime-picker
      v-model="birthday"
      type="date"
      :min-date="minDate"
      :max-date="maxDate"
      @confirm="onConfirm"
      @cancel="show = false"
    >
    </van-datetime-picker>
  </van-popup>
  <van-button class="button" round type="primary" size="large" @click="confirm"
    >确认</van-button
  >
</template>

<script setup lang="ts">
import axios from "axios";
import { upload } from "@h5/plugin/utils/upload";
import dataTool from "@h5/plugin/utils/date";
import * as dayjs from "dayjs";

import { STATIC_DIR } from "@h5/plugin/config";
import { ref, reactive } from "vue";
import { useUserStore } from "@h5/store/user";
import { useRouter } from "vue-router";

const store = useUserStore();
const router = useRouter();

const show = ref(false);
const loading = ref(false);
const birthday = ref(new Date(2000, 0, 1));
const minDate = new Date(1900, 0, 1);
const maxDate = new Date();
const staticDir = STATIC_DIR;

const user = reactive(store.auth);
if (!user.intro) user.intro = "我不想介绍自己";
if (!user.signature) user.signature = "太个性签名签不下";
user.birthday = dayjs(user.birthday).format(dataTool.formatYMD);
if (user.birthday !== dataTool.zeroTime)
  birthday.value = dayjs(user.birthday).toDate();

async function afterRead(file: any) {
  loading.value = true;
  user.avatarUrl = await upload(file.file);
  loading.value = false;
}
async function confirm() {
  user.birthday = dayjs(birthday.value).format(dataTool.formatYMD);
  const res = await axios.put(`/api/v1/user/${user.id}`, {
    id: user.id,
    details: {
      name: user.name,
      gender: user.gender,
      avatarUrl: user.avatarUrl,
      signature: user.signature,
      intro: user.intro,
      birthday: user.birthday,
    },
  });
  if (res.data.code == 0) store.auth = user;
  router.push("/me");
}
function onConfirm(value) {
  show.value = false;
  user.birthday = dayjs(birthday.value).format(dataTool.formatYMD);
}
</script>

<style scoped lang="less">
.van-row {
  padding: 1rem 0;
}
</style>
