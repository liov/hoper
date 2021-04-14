<template>
  <div class="wrapper">
    <van-uploader
      :deletable="false"
      :max-count="1"
      :preview-full-image="false"
      :after-read="afterRead"
    >
      <template #default>
        <van-image width="200" height="200" :src="url" fit="cover" />
      </template>
    </van-uploader>
    <van-loading v-if="loading" size="24px">上传中...</van-loading>
    <p>点击可更换头像</p>
    <div>
      <van-button class="button" plain type="success">关闭</van-button>
      <van-button class="button" plain type="primary" @click="confirm"
        >确认</van-button
      >
    </div>
  </div>
</template>

<script lang="ts">
import { Options, prop, Vue } from "vue-class-component";
import axios from "axios";
import { upload } from "@/plugin/utils/upload";

@Options({
  components: {},
})
export default class Edit extends Vue {
  user = null;
  url = "";
  loading = false;

  async created() {
    this.user = this.$store.state.auth;

    if (!this.user.introduction) this.user.introduction = "我不想介绍自己";
    if (!this.user.signature) this.user.signature = "太个性签名签不下";

    this.url = this.user.avatarUrl;
  }
  avatar() {
    this.show = true;
    this.url = this.user.avatarUrl;
  }
  async afterRead(file: any) {
    this.loading = true;
    this.url = await upload(file.file);
    this.loading = false;
  }
  async confirm() {
    await axios.post(`/api/v1/user`, {
      id: this.user.id,
      details: {
        avatarUrl: this.url,
      },
    });
  }
}
</script>

<style scoped></style>
