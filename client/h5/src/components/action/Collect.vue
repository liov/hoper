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
      <span class="button" @click="show = false">取消</span>
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

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
import Action from "@/components/action/Action.vue";
import axios from "axios";
import { upload } from "@/plugin/utils/upload";

@Options({ components: { Action } })
export default class AddCollect extends Vue {
  show = false;
  addShow = false;
  favs: any[] = [];
  refId = 0;
  type = 0;
  collects = [];
  collectsRef = null;
  loading = false;
  uploader = [];
  title = "";
  setCollect(param) {
    this.type = param.type;
    this.refId = param.refId;
    this.collects = param.collects;
    this.collectsRef = param.collects;
  }
  async created() {
    const res = await axios.get("/api/v1/content/tinyFav/0");
    console.log(res);
    this.favs = res.data.details.list;
  }
  async onCollect() {
    await axios.post("/api/v1/action/collect", {
      type: this.type,
      refId: this.refId,
      favIds: this.collects,
    });
    this.$toast.success("收藏成功");
    this.collectsRef.push(...this.collects);
    this.show = false;
  }
  async afterRead(file: any) {
    this.loading = true;
    file.url = await upload(file.file);
    this.loading = false;
  }
  async onConfirm() {
    if (this.title.trimStart().trimEnd().length === 0) {
      this.$toast.fail("内容为空");
      return;
    }
    const fav = {
      title: this.title,
      cover: this.uploader.length > 0 ? this.uploader[0].url : "",
    };
    const res = await axios.post("/api/v1/content/fav", fav);
    fav.id = res.data.details.id;
    fav.userId = this.$store.state.user.auth.id;
    this.favs.push(fav);
    this.$toast.success("新建成功");
    this.title = "";
  }
}
</script>

<style scoped lang="less">
.button {
  margin: 1rem 2rem;
}
</style>
