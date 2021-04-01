<template>
  <div>
    <van-list
      v-model:loading="loading"
      :finished="finished"
      finished-text="没有更多了"
      @load="onLoad"
    >
      <van-cell v-for="item in list">
        <template #default>
          <van-skeleton title avatar round :row="3" :loading="loading">
            <div class="moment" v-if="show">
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
                  v-for="img in item.images.split(',')"
                  :src="img"
                  lazy-load
                  class="img"
                />
              </lazy-component>
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
import { ObjMap } from "@/plugin/utils/user";

@Options({
  components: {},
})
export default class MomentList extends Vue {
  active = 0;
  loading = false;
  finished = false;
  show = false;
  pageNo = 1;
  pageSize = 10;
  userM = new ObjMap();
  list = Array.from(new Array(this.pageSize), (v, i) => {
    return { id: i };
  });

  //mounted() {}

  async onLoad() {
    // 异步更新数据
    const res = await axios.get(
      `/api/v1/moment?pageNo=${this.pageNo}&pageSize=${this.pageSize}`
    );
    const data = res.data.details;
    if (this.pageNo == 1) {
      this.list = data.list;
      this.userM.appendMap(data.users);
    } else {
      this.list = this.list.concat(data.list);
      this.userM.appendMap(data.users);
    }
    this.loading = false;
    this.show = true;
    this.pageNo++;
    if (data.list.length < this.pageSize) this.finished = true;
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
}
</style>
