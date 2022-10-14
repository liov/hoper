<template>

</template>


<script setup lang="ts">
  import { Toast } from "@/plugins/compatible/toast";
  import { defineAsyncComponent, ref } from "vue";
  const active = ref(2);

  function onChange(index: number) {
  Toast.text(index.toString());
}
</script>

<style scoped lang="scss">
i {
  margin: 0.7rem;
}
</style>
