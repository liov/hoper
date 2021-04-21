<template>
  <Moment></Moment>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import Moment from "@/components/moment/Moment.vue";
@Options({
  components: {
    Moment,
  },
})
export default class MomentDetail extends Vue {
  active = 0;

  onChange(index: number) {
    this.$toast({ message: index });
  }
}
</script>

<style scoped></style>
