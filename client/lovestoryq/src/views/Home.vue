<template>
  <div class="home">
    <van-tabs v-model:active="active" @change="onChange" swipeable title-active-color="#1989fa" animated line-height="0">
      <template #nav-left>
        <div style="width:30%">
        </div>
      </template>
      <template #nav-right>
        <div style="width:30%">
        </div>
      </template>
      <van-tab title='关注'>
        内容1
      </van-tab>
      <van-tab title='推荐'>
        内容2
      </van-tab>
      <van-tab title='最新'>
        内容3
      </van-tab>
    </van-tabs>
    <van-row type="flex">
      <van-col span="6">
        span: 6
      </van-col>
      <van-col span="12">
      </van-col>
      <van-col span="6">span: 6</van-col>
    </van-row>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import HelloWorld from "@/components/HelloWorld.vue"; // @ is an alias to /src
import { Toast, Notify } from "vant";

@Options({
  components: {
    HelloWorld
  }
})
export default class Home extends Vue {
  active = 0;
  onChange(index:number){
    Toast({message:index})
  }
}
</script>
