<template>
  <div></div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import axios from "axios";
import store from "@/store/index";

@Options({
  components: {},
})
export default class Active extends Vue {
  mounted() {
    axios
      .get(
        `/api/v1/user/active/${this.$route.params.id}/${this.$route.params.secret}`
      )
      .then((res) => {
        if (!res.data.code || res.data.code === 0) {
          this.$toast.success(res.data.message);
          this.$router.push({ path: "/" });
        } else this.$toast.fail(res.data.message);
      });
  }
}
</script>

<style scoped></style>
