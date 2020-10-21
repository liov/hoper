<template>
  <div>
    <van-field
      v-model="message"
      rows="3"
      autosize
      label-width="0"
      type="textarea"
      maxlength="567"
      size="large"
      show-word-limit
    />
    <van-field name="uploader" label-width="0">
      <template #input>
        <van-uploader
          v-model="uploader"
          :max-count="9"
          :max-size="500 * 1024"
          @oversize="onOversize"
          :after-read="afterRead"
        />
      </template>
    </van-field>
    <van-field
      readonly
      clickable
      label="权限"
      :value="permission"
      placeholder="选择权限"
      @click="showPicker = true"
    />
    <van-popup v-model:show="showPicker" round position="bottom">
      <van-picker
        :columns="columns"
        @cancel="showPicker = false"
        @confirm="onConfirm"
      />
    </van-popup>
    <div style="margin: 16px;">
      <van-button
        round
        block
        type="primary"
        native-type="submit"
        @click="submit"
      >
        提交
      </van-button>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import axios from "axios";
@Options({})
export default class MomentAdd extends Vue {
  message = "";
  permission = "";
  columns = ["全部", "自己可见", "陌生人可见"];
  showPicker = false;
  uploader = [{ url: "https://img.yzcdn.cn/vant/leaf.jpg" }];

  onOversize(file: File) {
    console.log(file);
    this.$toast("文件大小不能超过 500kb");
  }
  afterRead(file: File) {
    // 此时可以自行将文件上传至服务器
    console.log(file);
  }
  async submit() {
    const res = await axios.post(`/api/moment`, {
      image_url: "",
      mood_name: "",
      tags: [],
      permission: 0,
      content: this.message
    });
    if (res.data.code === 200) {
      await this.$router.push({ path: "/moment" });
    }
  }
  onConfirm(value: string) {
    console.log(value)
    this.permission = value;
    this.showPicker = false;
  }
}
</script>

<style scoped></style>
