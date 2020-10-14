<template>
  <div>

  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import {Toast} from "vant";

@Options({
  components: {
  }
})
export default class MomentList extends Vue {
  active = 0;
  onChange(index:number){
    Toast({message:index})
  }
}
</script>

<style scoped>

</style>