<template>
  <div>
    <van-field
      v-model="message"
      rows="2"
      autosize
      type="textarea"
      maxlength="50"
      show-word-limit
    />
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
@Options({})
export default class Moment extends Vue {
  message = "";
}
</script>

<style scoped></style>
