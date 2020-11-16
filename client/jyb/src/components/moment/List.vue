<template>
  <view>
    <at-virtual-scroll
      bench="5"
      :items="list"
      item-height="200"
      @onReachBottom="onLoad"
    >
      <template #default="{ index, item }">
        <AtCard
          note='小Tips'
          :extra='this.$date2s(item.created_at)'
          :extraStyle="{left:'100px',textOverflow: 'clip',maxWidth: '10rem'}"
          :title='item.user.name'
          :thumb='item.user.avatar_url'
        >
          <AtNoticebar marquee>{{item.content}}</AtNoticebar>
        </AtCard>
      </template>
    </at-virtual-scroll>
  </view>
</template>

<script lang="ts">
import {Options, Vue} from "vue-class-component";
import axios from "axios";
import { AtVirtualScroll,AtCard,AtNoticebar  } from 'taro-ui-vue3'
@Options({
  components: {AtVirtualScroll,AtCard,AtNoticebar }
})
export default class MomentList extends Vue {
  active = 0;
  list: any[] = [];
  loading = false;
  finished = false;
  show = false;
  pageNo = 0;
  pageSize = 10;

  mounted() {
    this.onLoad()
  }

  async onLoad() {
    // 异步更新数据
    const res = await axios.get(
      `/api/moment?pageNo=${this.pageNo}&pageSize=${this.pageSize}`
    );
    if (this.pageNo == 0) {
      this.list = res.data.data;
    } else this.list = this.list.concat(res.data.data);
    this.loading = false;
    this.show = true;
    this.pageNo++;
    if (res.data.data.length < this.pageSize) {
      this.finished = true;
    }
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
