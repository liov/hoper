<template>
  <div class="home">
    <van-row>
      <van-col span="8"
        ><van-image round width="10rem" height="10rem" :src="user.avatarURL"
      /></van-col>
      <van-col span="8">
        <van-row>
          <van-col span="8">span: 8</van-col>
          <van-col span="8">span: 8</van-col>
        </van-row>
        <van-row>
          <van-col span="8">span: 8</van-col>
          <van-col span="8">span: 8</van-col>
          <van-col span="8">span: 8</van-col>
        </van-row></van-col
      >
    </van-row>
    <van-row>
      <van-col span="8">span: 8</van-col>
      <van-col span="8">span: 8</van-col>
      <van-col span="8">span: 8</van-col>
    </van-row>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import axios from "axios";
@Options({
  components: {},
})
export default class Home extends Vue {
  user = {};
  async created() {
    const res = await axios.get(`/api/v1/user/id/${this.$store.state.auth.id}`);
    this.user = res.data.details.user;
  }
  mounted() {
    console.log(this.user);
  }
}
</script>

<style scoped lang="less">
i {
  margin: 0.7rem;
}
</style>
