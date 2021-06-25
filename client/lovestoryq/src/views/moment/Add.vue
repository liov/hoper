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
      v-model="permissionVal"
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
    <div style="margin: 16px">
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
import { upload } from "@/plugin/utils/upload";

@Options({ components: {} })
export default class MomentAdd extends Vue {
  message = "";
  permission = 0;
  permissionVal = "全部";
  columns = ["全部", "自己可见", "陌生人可见"];
  showPicker = false;
  uploader: any = [];

  onOversize(file: any) {
    if (file.file.size > 5 * 1e5) this.$toast("文件大小不能超过 500kB");
  }
  async afterRead(file: any) {
    file.url = await upload(file.file);
  }
  async submit() {
    let images = "";
    for (const up of this.uploader) {
      images += up.url + ",";
    }
    if (images !== "") images = images.slice(0, images.length - 1);
    const res = await axios.post(`/api/v1/moment`, {
      mood: "",
      tags: [],
      permission: this.permission,
      anonymous: 2,
      content: this.message,
      images: images,
    });
    if (res.data.code === 0) {
      await this.$router.push({ path: "/" });
    }
  }
  onConfirm(value: string, index: number) {
    this.permission = index;
    this.permissionVal = value;
    this.showPicker = false;
  }
}
</script>

<style scoped></style>
