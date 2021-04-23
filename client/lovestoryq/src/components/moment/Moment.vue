<template>
  <div class="moment">
    <div class="auth">
      <img class="avatar" :src="user.avatarUrl" />
      <span class="name">{{ user.name }}</span>
      <span class="time">{{ $date2s(moment.createdAt) }}</span>
    </div>
    <div class="content" @click="detail">
      <van-field
        v-model="moment.content"
        rows="1"
        :autosize="maxHeight ? { maxHeight } : true"
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
    <lazy-component class="imgs" v-if:="moment.images">
      <van-image
        width="100"
        height="100"
        v-for="(img, idx) in images"
        :src="img"
        lazy-load
        class="img"
        @click="preview(idx)"
      />
    </lazy-component>
    <Action :content="moment" :type="1"> </Action>
  </div>
</template>

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
import { ImagePreview } from "vant";
import Action from "@/components/action/Action.vue";
import { jump } from "@/router/utils";
class Props {
  moment = prop<any>({ default: {} });
  user = prop<any>({});
  maxHeight = prop<number>({});
}
@Options({ components: { Action } })
export default class Moment extends Vue.with(Props) {
  images = [];
  created() {
    this.images = this.moment.images?.split(",");
  }
  preview(idx: number) {
    ImagePreview({
      images: this.images,
      startPosition: idx,
      closeable: true,
    });
  }
  detail() {
    jump(this.$route.path, 1, this.moment);
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
