<template>
  <van-action-sheet
    :show="show"
    :actions="actions"
    cancel-text="取消"
    close-on-click-action
    @click="show = !show"
    teleport="#app"
  >
  </van-action-sheet>
  <van-dialog
    v-model:show="report.show"
    title="举报"
    show-cancel-button
    teleport="#app"
    @confirm="onReport"
  >
    <van-field name="radio">
      <template #input>
        <van-radio-group
          v-model="report.checked"
          direction="horizontal"
          @change="remark"
        >
          <van-radio
            v-for="item in report.options"
            :name="item.id"
            shape="square"
            >{{ item.name }}</van-radio
          >
        </van-radio-group>
      </template>
    </van-field>
    <van-field
      v-show="report.field"
      v-model="report.message"
      rows="1"
      autosize
      label="备注"
      type="textarea"
      placeholder="请输入举报内容"
    />
  </van-dialog>
  <van-share-sheet
    :show="share.show"
    title="立即分享给好友"
    @cancel="share.show = !share.show"
    :options="share.options"
    teleport="#app"
  />
  <AddCollect ref="addCollect"></AddCollect>
</template>

<script setup lang="ts">
import axios from "axios";
import { onMounted, onUnmounted, reactive, ref } from "vue";
import emitter from "@/plugin/emitter";
import AddCollect from "@/components/action/Collect.vue";
import { Dialog } from "vant";

const VanDialog = Dialog.Component;
let type = 0;
let refId = 0;
const show = ref(false);
const addCollect: any = ref(null);

const actions = [
  { name: "分享", callback: () => (share.show = !share.show) },
  { name: "不喜欢" },
  { name: "举报", callback: () => (report.show = !report.show) },
  { name: "删除", color: "#D91E46" },
];
const report = reactive({
  show: false,
  checked: false,
  field: false,
  message: "",
  options: [
    { id: 1, name: "色情暴力" },
    { id: 2, name: "侮辱谩骂" },
    { id: 3, name: "涉及政治" },
    { id: 255, name: "其他原因" },
  ],
});

const share = reactive({
  show: false,
  options: [
    [
      { name: "微信", icon: "wechat" },
      { name: "朋友圈", icon: "wechat-moments" },
      { name: "微博", icon: "weibo" },
      { name: "QQ", icon: "qq" },
    ],
    [
      { name: "复制链接", icon: "link" },
      { name: "分享海报", icon: "poster" },
      { name: "二维码", icon: "qrcode" },
      { name: "小程序码", icon: "weapp-qrcode" },
    ],
  ],
});
const comment = reactive({
  show: false,
  replyId: 0,
});
onMounted(() => {
  emitter.on("more-show", (param: any) => {
    type = param.type;
    refId = param.refId;
    show.value = !show.value; //箭头函数内部不会产生新的this，这边如果不用=>,this指代Event
  });
  emitter.on("fav-show", (param) => {
    addCollect.value.setCollect(param);
    addCollect.value.show = true;
  });
  console.log(emitter.all);
});
onUnmounted(() => {
  emitter.all.delete("more-show");
  emitter.all.delete("fav-show");
});
function remark(name: string) {
  report.field = name === "255";
}
async function onReport() {
  await axios.post("/api/v1/action/report", {
    type,
    refId,
    remark: report.message,
  });
}
</script>

<style scoped></style>
