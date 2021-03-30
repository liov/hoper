<template>
  <div>
    <van-list
      v-model:loading="loading"
      :finished="finished"
      finished-text="没有更多了"
      @load="onLoad"
    >
      <van-cell v-for="(item, index) in list" :key="index">
        <template #default>
          <van-skeleton title avatar :row="3" :loading="loading">
            <div class="moment" v-if="show">
              <img :src="this.userM.get(item.userId).avatarUrl" />
              <div class="content">
                <span>{{ this.userM.get(item.userId).name }}</span>
                <span>{{ $date2s(this.userM.get(item.userId).created_at) }}</span>
                <div class="van-multi-ellipsis--l3">
                  {{ item.content }}
                </div>
              </div>
            </div>
          </van-skeleton>
        </template>
      </van-cell>
    </van-list>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import axios from "axios";
import {UserBaseInfo,appendUserMap} from "@/plugin/utils/user";

@Options({
  components: {}
})
export default class MomentList extends Vue {
  active = 0;
  list: any[] = [];
  loading = false;
  finished = false;
  show = false;
  pageNo = 1;
  pageSize = 10;
  userM = new Map<number,UserBaseInfo>();

  mounted() {
    if (this.list.length == 0) {
      const list: any[] = [];
      for (let i = 0; i < 5; i++) {
        list.push({ id: i });
      }
      this.list = list;
    }
  }

  async onLoad() {
    // 异步更新数据
    const res = await axios.get(
      `/api/v1/moment?pageNo=${this.pageNo}&pageSize=${this.pageSize}`
    );
    const data = res.data.details;
    if (this.pageNo == 1) {
      this.list = data.list;
      appendUserMap(this.userM,data.users)
    } else {
      this.list = this.list.concat(data.list);
      appendUserMap(this.userM,data.users)
    }
    this.loading = false;
    this.show = true;
    this.pageNo++;
    if (data.list.length < this.pageSize) this.finished = true;
  }
}
</script>

<style scoped lang="scss">
.moment {
  display: flex;
  padding: 0 16px;

  .content {
    padding-top: 6px;

    h3 {
      margin: 0;
      font-size: 18px;
      line-height: 20px;
    }

    .van-multi-ellipsis--l3 {
      margin: 13px 0 0;
      font-size: 14px;
      line-height: 20px;
    }
  }

  img {
    flex-shrink: 0;
    width: 32px;
    height: 32px;
    margin-right: 16px;
  }
}
</style>
