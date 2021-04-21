<template>
  <van-popup
    v-model:show="show"
    position="bottom"
    teleport="#app"
    :overlay="false"
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
</template>

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
import Action from "@/components/action/Action.vue";
import axios from "axios";

@Options({ components: { Action } })
export default class AddCollect extends Vue {
  show = false;
  message = "";
  favs = [];
  refId = 0;
  type = 0;
  collects = [];
  setCollect(param) {
    console.log(param);
    this.type = param.type;
    this.refId = param.refId;
    this.collects = param.collects ? param.collects : [];
  }
  async created() {
    const res = await axios.get("/api/v1/content/tinyFav/0");
    this.favs = res.data.details.list;
  }
  async onCollect() {
    await axios.post("/api/v1/action/collect", {
      type: this.type,
      refId: this.refId,
      favIds: this.collects,
    });
    this.$toast.success("收藏成功");
    this.show = false;
  }
}
</script>

<style scoped lang="less">
.button {
  margin: 1rem 2rem;
}
</style>
