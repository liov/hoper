<template>
  <div class="moment">
    <Moment v-if="show" :moment="moment" :user="user"></Moment>
  </div>
  <CommentList :type="1" :refId="$route.params.id"></CommentList>
  <AddComment
    v-if="show"
    ref="addComment"
    :comment="{ type: 1, refId: moment.id, recvId: user.id }"
  ></AddComment>
  <div class="placeholder"></div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import Moment from "@/components/moment/Moment.vue";
import ActionMore from "@/components/action/More.vue";
import CommentList from "@/components/comment/List.vue";
import AddComment from "@/components/comment/Add.vue";
import axios from "axios";
@Options({
  components: {
    Moment,
    ActionMore,
    CommentList,
    AddComment,
  },
})
export default class MomentDetail extends Vue {
  active = 0;
  moment = null;
  user = null;
  show = false;
  async created() {
    this.show = false;
    this.moment = this.$store.state.content.moment;
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
.placeholder {
  height: 100px;
}
</style>
