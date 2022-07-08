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

<script setup lang="ts">
import { ref } from "vue";
import Moment from "@/components/moment/Moment.vue";
import CommentList from "@/components/comment/List.vue";
import AddComment from "@/components/comment/Add.vue";
import axios from "axios";
import { useContentStore } from "@/store/content";
import { useRoute } from "vue-router";
import { useUserStore } from "@/store/user";

const store = useContentStore();
const userStore = useUserStore();
const route = useRoute();
const active = ref(0);

const moment = ref(store.moment);

if (!moment.value) {
  const res = await axios.get(`/api/v1/moment/${route.params.id}`);
  moment.value = res.data.details;
  store.moment = moment.value;
  userStore.appendUsers(moment.value!.users);
}
const user = getUser(moment.value!.userId);
const show = true;

function getUser(id: number) {
  return userStore.getUser(id);
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
