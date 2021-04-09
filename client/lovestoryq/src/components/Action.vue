<template>
  <van-row>
    <van-col span="6" class="action" @click="show.moreShow = !show.moreShow"
      ><van-icon name="more-o"
    /></van-col>
    <van-col span="6" class="action"
      ><van-icon
        :name="content.collect ? 'star' : 'star-o'"
        :color="content.collect ? '#F6DF02' : ''"
    /></van-col>
    <van-col span="6" class="action"><van-icon name="comment-o" /></van-col>
    <van-col span="6" class="action"
      ><van-icon
        :name="content.likeId > 0 ? 'like' : 'like-o'"
        :color="content.likeId > 0 ? '#D91E46' : ''"
        @click="like(index)"
    /></van-col>
  </van-row>
  <van-action-sheet
    :show="show.moreShow"
    :actions="report.actions"
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
  >
    <van-field name="radio">
      <template #input>
        <van-radio-group
          v-model="report.checked"
          direction="horizontal"
          @change="remark"
        >
          <van-radio name="1" shape="square">色情暴力</van-radio>
          <van-radio name="2" shape="square">侮辱谩骂</van-radio>
          <van-radio name="3" shape="square">政治政策</van-radio>
          <van-radio name="255" shape="square">其他原因</van-radio>
        </van-radio-group>
      </template>
    </van-field>
    <van-field
      v-if="report.field"
      v-model="report.message"
      rows="1"
      autosize
      label="备注"
      type="textarea"
      placeholder="请输入举报内容"
    />
  </van-dialog>
</template>

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
import axios from "axios";
import { reactive } from "vue";
class Props {
  content = prop<any>({ default: {} });
  readonly type: number = prop<number>();
}
@Options({})
export default class Action extends Vue.with(Props) {
  show = reactive({
    moreShow: false,
    shareShow: false,
  });

  report = reactive({
    checked: false,
    show: false,
    actions: [
      { name: "不喜欢" },
      {
        name: "举报",
        callback: () => (this.report.show = !this.report.show),
      },
      {
        name: "删除",
        color: "#D91E46",
      },
    ],
    field: false,
    message: "",
  });
  remark(name: string) {
    console.log(name);
    if (name === "255") {
      this.report.field = true;
    }
  }
  async like() {
    const api = `/api/v1/action/like`;
    const id = this.content.id;
    const likeId = this.content.likeId;
    if (likeId > 0) {
      await axios.delete(api, { data: { id: likeId } });
      this.content.likeId = 0;
    } else {
      const res = await axios.post(api, {
        refId: id,
        type: this.type,
        action: 2,
      });
      this.content.likeId = res.data.details.id;
    }
  }
}
</script>

<style scoped>
.action {
  text-align: center;
}
</style>
