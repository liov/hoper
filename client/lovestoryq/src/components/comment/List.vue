<template>
  <div>
    <van-list
      v-model:loading="loading"
      :finished="finished"
      finished-text="没有更多了"
      @load="onLoad"
    >
      <van-cell v-for="item in list" :key="item.id">
        <template #default>
          <van-skeleton title avatar round :row="3" :loading="loading">
            <Comment
              v-if="show"
              :comment="item"
              :user="user(item.userId)"
            ></Comment>
          </van-skeleton>
        </template>
      </van-cell>
    </van-list>
    <ActionMore></ActionMore>
  </div>
</template>

<script lang="ts">
import { Options, prop, Vue } from "vue-class-component";
import axios from "axios";
import { reactive, ref } from "vue";
import Comment from "@/components/comment/Comment.vue";
import ActionMore from "@/components/action/More.vue";

class Props {
  type = prop<number>({});
  refId = prop<number>({});
  rootId = prop<number>({ default: 0 });
}
@Options({
  components: { Comment, ActionMore },
})
export default class CommentList extends Vue.with(Props) {
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
      `/api/v1/action/comment?type=${this.type}&refId=${this.refId}&rootId=${this.rootId}&pageNo=${this.pageNo}&pageSize=${this.pageSize}`
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
    this.$store.state.content.commentCache.set(this.rootId, this.list);
    this.show = true;
    this.$store.commit("appendUsers", data.users);
    this.pageNo++;
    if (data.list.length < this.pageSize) this.finished = true;
  }
}
</script>

<style scoped lang="less"></style>
