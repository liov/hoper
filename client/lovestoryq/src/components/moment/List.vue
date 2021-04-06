<template>
  <div>
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
            v-model="checked"
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
    <van-pull-refresh
      v-model="state.loading"
      :success-text="successText"
      @refresh="onRefresh"
    >
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-cell v-for="(item, index) in list">
          <template #default>
            <van-skeleton title avatar round :row="3" :loading="loading">
              <div class="moment" v-if="show.listShow">
                <div class="auth">
                  <img class="avatar" :src="userM.get(item.userId).avatarUrl" />
                  <span class="name">{{ userM.get(item.userId).name }}</span>
                  <span class="time">{{ $date2s(item.createdAt) }}</span>
                </div>
                <div class="content">
                  <van-field
                    v-model="item.content"
                    rows="1"
                    :autosize="{ maxHeight: 200 }"
                    readonly
                    type="textarea"
                  >
                    <template #extra>
                      <div class="arrow">
                        <van-icon name="arrow-down" />
                      </div>
                    </template>
                  </van-field>
                </div>
                <lazy-component class="imgs" v-if:="item.images">
                  <van-image
                    width="100"
                    height="100"
                    v-for="(img, idx) in item.images.split(',')"
                    :src="img"
                    lazy-load
                    class="img"
                    @click="preview(idx, item.images)"
                  />
                </lazy-component>
              </div>

              <van-row>
                <van-col
                  span="6"
                  class="action"
                  @click="show.moreShow = !show.moreShow"
                  ><van-icon name="more-o"
                /></van-col>
                <van-col span="6" class="action"
                  ><van-icon
                    :name="item.collect ? 'star' : 'star-o'"
                    :color="item.collect ? '#F6DF02' : ''"
                /></van-col>
                <van-col span="6" class="action"
                  ><van-icon name="comment-o"
                /></van-col>
                <van-col span="6" class="action"
                  ><van-icon
                    :name="item.likeId > 0 ? 'like' : 'like-o'"
                    :color="item.likeId > 0 ? '#D91E46' : ''"
                    @click="like(index)"
                /></van-col>
              </van-row>
            </van-skeleton>
          </template>
        </van-cell>
      </van-list>
    </van-pull-refresh>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import axios from "axios";
import { ObjMap } from "@/plugin/utils/user";
import { ImagePreview } from "vant";
import { reactive, ref } from "vue";

@Options({
  components: {},
})
export default class MomentList extends Vue {
  active = 0;
  loading = false;
  finished = false;
  pageNo = 1;
  pageSize = 10;
  userM = new ObjMap();
  list = Array.from(new Array(this.pageSize), (v, i) => {
    return { id: i };
  });
  successText = "刷新成功";
  state = reactive({
    loading: false,
  });
  show = reactive({
    listShow: ref(false),
    moreShow: ref(false),
    shareShow: ref(false),
  });
  checked = ref(false);

  report = {
    show: ref(false),
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
    field: ref(false),
    message: "",
  };
  //mounted() {}

  async onLoad() {
    this.finished = false;
    // 异步更新数据
    const res = await axios.get(
      `/api/v1/moment?pageNo=${this.pageNo}&pageSize=${this.pageSize}`
    );
    if (res.data.code !== 0) {
      this.$toast.fail(res.data.message);
      this.finished = true;
    }
    const data = res.data.details;
    if (this.pageNo == 1) {
      this.list = data.list;
      this.userM.appendMap(data.users);
    } else {
      this.list = this.list.concat(data.list);
      this.userM.appendMap(data.users);
    }
    this.loading = false;
    this.show.listShow = true;
    this.pageNo++;
    if (data.list.length < this.pageSize) this.finished = true;
  }
  preview(idx: number, images: string) {
    ImagePreview({
      images: images.split(","),
      startPosition: idx,
      closeable: true,
    });
  }
  onRefresh = () => {
    this.pageNo = 1;
    this.onLoad().catch(() => {
      this.successText = "刷新失败";
    });
    this.state.loading = false;
  };
  remark(name: string) {
    console.log(name);
    if (name === "255") {
      this.report.field = true;
    }
  }
  async like(idx: number) {
    console.log(this.list[idx]);
    const api = `/api/v1/action/like`;
    const id = this.list[idx].id;
    const likeId = this.list[idx].likeId;
    if (likeId > 0) {
      await axios.delete(api, { data: { id: likeId } });
      this.list[idx].likeId = 0;
    } else {
      this.list[idx].likeId = await axios.post(api, {
        refId: id,
        type: 1,
        action: 2,
      });
    }
  }
}
</script>

<style scoped lang="less">
.moment {
  @20px: 20px;
  @avatar: 30px;
  .name {
    left: 60px;
    position: absolute;
  }

  .time {
    position: absolute;
    right: @20px;
  }
  .content {
    width: 100%;
    h3 {
      margin: 0;
      font-size: 18px;
      line-height: 20px;
    }

    .arrow {
      position: absolute;
      bottom: 16px;
      right: 0;
    }

    .van-multi-ellipsis--l3 {
      margin: 13px 0 0;
      font-size: 14px;
      line-height: 20px;
    }
  }

  .avatar {
    flex-shrink: 0;
    width: @avatar;
    height: @avatar;
    border-radius: 40px;
    position: relative;
    margin: 0 16px;
  }
  .imgs {
    padding: 0 11px;
  }
  .img {
    margin: 5px 5px;
  }
  .action {
    text-align: center;
  }
}
</style>
