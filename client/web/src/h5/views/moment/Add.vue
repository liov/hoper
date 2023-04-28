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
          @click="notH5Upload"
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

<script setup lang="ts">
import axios from "axios";
import { upload } from "@h5/plugin/utils/upload";
import { showToast } from "vant";
import { useRouter } from "vue-router";
import { reactive, ref } from "vue";
import { useGlobalStore } from "@h5/store/global";
import { Platform } from "@h5/model/const";

const globalState = useGlobalStore();

const router = useRouter();

const message = ref("");
const permission = ref(0);
const permissionVal = ref("全部");
const columns = [
  { text: "全部" },
  { text: "自己可见" },
  { text: "陌生人可见" },
];
const showPicker = ref(false);
const uploader: any = reactive([]);

function onOversize(file: any) {
  if (file.file.size > 5 * 1e5) showToast("文件大小不能超过 500kB");
}
async function afterRead(file: any) {
  file.url = await upload(file.file);
}
async function submit() {
  let images = "";
  for (const up of uploader) {
    images += up.url + ",";
  }
  if (images !== "") images = images.slice(0, images.length - 1);
  const res = await axios.post(`/api/v1/moment`, {
    mood: "",
    tags: [],
    permission: permission,
    anonymous: 2,
    content: message,
    images: images,
  });
  if (res.data.code === 0) {
    await router.push({ path: "/" });
  }
}
function onConfirm(value: string, index: number) {
  permission.value = index;
  permissionVal.value = value;
  showPicker.value = false;
}

function notH5Upload() {
  if (globalState.platform == Platform.App) {
    window.Flutter.postMessage(
      JSON.stringify({ method: "pickPhoto", params: [] })
    );
  }
}
</script>

<style scoped></style>
