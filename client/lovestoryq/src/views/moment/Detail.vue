<template>
  <div class="moment">
    <Moment v-if="show" :moment="moment" :user="user"></Moment>
    <ActionMore></ActionMore>
  </div>
  <CommentList :type="1" :refId="$route.params.id"></CommentList>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import Moment from "@/components/moment/Moment.vue";
import ActionMore from "@/components/action/More.vue";
import CommentList from "@/components/comment/List.vue";
import axios from "axios";
@Options({
  components: {
    Moment,
    ActionMore,
    CommentList,
  },
})
export default class MomentDetail extends Vue {
  active = 0;
  moment = null;
  user = null;
  show = false;
  async created() {
    this.show = false;
    this.moment = this.$store.state.moment.moment;
    if (!this.moment) {
      const res = await axios.get(`/api/v1/moment/${this.$route.params.id}`);
      this.moment = res.data.details;
      this.$store.commit("setMoment", this.moment);
      this.$store.commit("appendUsers", this.moment.users);
    }
    this.user = this.getUser(this.moment.userId);
    this.show = true;
  }
  getUser(id: number) {
    return this.$store.getters.getUser(id);
  }
}
</script>

<style scoped lang="less">
.moment {
  text-align: left;
}
</style>
