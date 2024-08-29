<template>
  <div></div>
</template>

<script setup lang="ts">
import axios from "axios";
import { onMounted } from "vue";
import { useUserStore } from "@/mixin/store/user";
import { useRouter, useRoute } from "vue-router";
import { Toast } from "vant";

const store = useUserStore();
const router = useRouter();
const route = useRoute();

onMounted(() => {
  axios
    .get(`/api/v1/user/active/${route.params.id}/${route.params.secret}`)
    .then((res) => {
      if (!res.data.code || res.data.code === 0) {
        Toast.success(res.data.msg);
        router.push({ path: "/" });
      } else Toast.fail(res.data.msg);
    });
});
</script>

<style scoped></style>
