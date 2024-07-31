<template>
  <van-popup
    v-model:show="show"
    position="bottom"
    teleport="#app"
    :overlay="false"
  >
    <van-icon name="add-o" size="16" @click="addShow = true"
      >新建收藏夹</van-icon
    >
    <van-checkbox-group v-model="collects">
      <van-checkbox v-for="item in favs" :name="item.id">{{
        item.title
      }}</van-checkbox>
    </van-checkbox-group>
    <div class="button">
      <span class="button" @click="show = !show">取消</span>
      <span class="button" style="color: red" @click="onCollect">确认</span>
    </div>
  </van-popup>
  <van-popup v-model:show="addShow" teleport="#app" position="bottom">
    <van-field
      v-model="title"
      placeholder="请输入收藏夹名称"
      :rules="[{ required: true, message: '输入内容为空' }]"
      ref="commentRef"
    >
      <template #button>
        <div class="button">
          <van-uploader
            v-model="uploader"
            :max-count="1"
            :after-read="afterRead"
          />
          <van-loading v-if="loading" size="24px">上传中...</van-loading>
          <van-button size="small" type="primary" @click="onConfirm"
            >确认</van-button
          >
        </div>
      </template>
    </van-field>
  </van-popup>
</template>

<script setup lang="ts">
import Action from "@h5/components/action/Action.vue";
import axios from "axios";
import { upload } from "@h5/plugin/utils/upload";
import { reactive, ref } from "vue";
import { Toast } from "vant";
import type { Ref } from "vue";
import { useUserStore } from "@h5/store/user";

const userStore = useUserStore();

const show = ref(false);
const addShow = ref(false);

let refId = 0;
let type = 0;
const collects: Ref<any[]> = ref([]);

const loading = ref(false);
const uploader: Ref<any[]> = ref([]);
const title = ref("");
function setCollect(param) {
  type = param.type;
  refId = param.refId;
  collects.value = param.collects;
}

defineExpose({
  show,
  addShow,
  loading,
  uploader,
  title,
  onCollect,
  afterRead,
  onConfirm,
  setCollect,
});

const res = await axios.get("/api/v1/content/tinyFav/0");
console.log(res);
const favs: UnwrapNestedRefs<any[]> = res.data.data.list
  ? reactive(res.data.data.list)
  : reactive([]);

async function onCollect() {
  await axios.post("/api/v1/action/collect", {
    type: type,
    refId: refId,
    favIds: collects,
  });
  Toast.success("收藏成功");
  show.value = false;
}
async function afterRead(file: any) {
  loading.value = true;
  file.url = await upload(file.file);
  loading.value = false;
}
async function onConfirm() {
  if (title.value.trimStart().trimEnd().length === 0) {
    Toast.fail("内容为空");
    return;
  }
  const fav: any = {
    title: title.value,
    cover: uploader.value.length > 0 ? uploader.value[0].url : "",
  };
  const res = await axios.post("/api/v1/content/fav", fav);
  fav.id = res.data.data.id;
  fav.userId = userStore.auth.id;
  favs.push(fav);
  Toast.success("新建成功");
  title.value = "";
}
</script>

<style scoped lang="less">
.button {
  margin: 1rem 2rem;
}
</style>
