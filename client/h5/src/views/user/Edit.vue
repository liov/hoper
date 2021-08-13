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
              :src="user.avatarUrl"
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
              <van-radio :name="1">男</van-radio>
              <van-radio :name="2">女</van-radio>
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

<script lang="ts">
import { Options, prop, Vue } from "vue-class-component";
import axios from "axios";
import { upload } from "@/plugin/utils/upload";
import dataTool from "@/plugin/utils/date";
import dayjs from "dayjs";
import store from "@/store";

@Options({
  components: {},
})
export default class Edit extends Vue {
  user = {};
  show = false;
  loading = false;
  birthday = new Date(2000, 0, 1);
  minDate = new Date(1900, 0, 1);
  maxDate = new Date();
  async created() {
    this.user = this.$store.state.user.auth;
    if (!this.user.intro) this.user.intro = "我不想介绍自己";
    if (!this.user.signature) this.user.signature = "太个性签名签不下";
    this.user.birthday = dayjs(this.user.birthday).format(dataTool.formatYMD);
    if (this.user.birthday !== dataTool.zeroTime)
      this.birthday = dayjs(this.user.birthday).toDate();
  }

  async afterRead(file: any) {
    this.loading = true;
    this.user.avatarUrl = await upload(file.file);
    this.loading = false;
  }
  async confirm() {
    this.user.birthday = this.birthday.format(dataTool.formatYMD);
    const res = await axios.put(`/api/v1/user/${this.user.id}`, {
      id: this.user.id,
      details: {
        name: this.user.name,
        gender: this.user.gender,
        avatarUrl: this.user.avatarUrl,
        signature: this.user.signature,
        intro: this.user.intro,
        birthday: this.user.birthday,
      },
    });
    if (res.data.code == 0) this.$store.commit("SET_AUTH", this.user);
  }
  onConfirm(value) {
    this.show = false;
    this.user.birthday = dayjs(this.birthday).format(dataTool.formatYMD);
  }
}
</script>

<style scoped lang="less">
.van-row {
  padding: 1rem 0;
}
</style>
