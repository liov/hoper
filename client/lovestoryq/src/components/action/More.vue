<template>
  <van-action-sheet
    :show="show.moreShow"
    :actions="actions"
    cancel-text="取消"
    close-on-click-action
    @click="show.moreShow = !show.moreShow"
    teleport="#app"
  >
  </van-action-sheet>
  <van-dialog
    :show="report.show"
    title="举报"
    show-cancel-button
    @cancel="report.show = !report.show"
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
  <CommentAdd ref="commentAdd"></CommentAdd>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import axios from "axios";
import { reactive } from "vue";
import emitter from "@/plugin/emitter";
import CommentAdd from "@/components/comment/Add.vue";
@Options({ components: { CommentAdd } })
export default class ActionMore extends Vue {
  type = 0;
  refId = 0;
  show = reactive({
    moreShow: false,
    favShow: false,
    commentShow: false,
  });

  actions = [
    { name: "分享", callback: () => (this.share.show = !this.share.show) },
    { name: "不喜欢" },
    { name: "举报", callback: () => (this.report.show = !this.report.show) },
    { name: "删除", color: "#D91E46" },
  ];
  report = reactive({
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

  share = reactive({
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
  comment = reactive({
    show: false,
    replyId: 0,
  });
  created() {
    emitter.on("more-show", (param) => {
      this.type = param.type;
      this.refId = param.refId;
      this.show.moreShow = !this.show.moreShow; //箭头函数内部不会产生新的this，这边如果不用=>,this指代Event
    });
    emitter.on("fav-show", (param) => {
      this.type = param.type;
      this.refId = param.refId;
      this.show.favShow = !this.show.favShow;
    });
    emitter.on("comment-show", (param) => {
      this.$refs.commentAdd.setComment(param);
      this.$refs.commentAdd.show = true;
    });
  }
  unmounted() {
    emitter.all.delete("more-show");
    emitter.all.delete("fav-show");
    emitter.all.delete("comment-show");
  }
  remark(name: string) {
    if (name === "255") this.report.field = true;
    else this.report.field = false;
  }
  async onReport() {
    await axios.post("/api/v1/action/report", {
      type: this.type,
      refId: this.refId,
      remark: this.report.message,
    });
  }
}
</script>

<style scoped></style>
