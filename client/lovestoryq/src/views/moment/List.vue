<template>
  <div>
    <van-pull-refresh
      v-model="pullDown.refreshing"
      :success-text="pullDown.successText"
      @refresh="onRefresh"
    >
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-cell v-for="item in list" :key="item.id">
          <template #default>
            <van-skeleton title avatar round :row="3" :loading="loading">
              <Moment
                v-if="show"
                :moment="item"
                :user="user(item.userId)"
                :maxHeight="200"
              ></Moment>
            </van-skeleton>
          </template>
        </van-cell>
      </van-list>
    </van-pull-refresh>
    <ActionMore></ActionMore>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import axios from "axios";
import { reactive, ref } from "vue";
import Moment from "@/components/moment/Moment.vue";
import ActionMore from "@/components/action/More.vue";

@Options({
  components: { Moment, ActionMore },
})
export default class MomentList extends Vue {
  active = 0;
  loading = false;
  finished = false;
  pageNo = 1;
  pageSize = 10;
  list = Array.from(new Array(this.pageSize), () => {
    return {};
  });

  pullDown = reactive({
    refreshing: false,
    successText: "刷新成功",
  });
  show = false;

  //mounted() {}
  user(id: number) {
    return this.$store.getters.getUser(id);
  }

  async onLoad() {
    this.finished = false;
    // 异步更新数据
    const res = await axios.get(
      `/api/v1/moment?pageNo=${this.pageNo}&pageSize=${this.pageSize}`
    );
    this.loading = false;
    const data = res.data.details;
    if (!data || !data.list) {
      this.finished = true;
      return;
    }
    if (this.pageNo == 1) {
      this.list = data.list;
    } else {
      this.list = this.list.concat(data.list);
    }
    this.$store.commit("appendUsers", data.users);
    this.show = true;
    this.pageNo++;
    if (data.list.length < this.pageSize) this.finished = true;
  }
  onRefresh = () => {
    this.pullDown.refreshing = true;
    this.pageNo = 1;
    this.onLoad().catch(() => {
      this.pullDown.successText = "刷新失败";
    });
    this.pullDown.refreshing = false;
  };
}
</script>

<style scoped lang="less"></style>
