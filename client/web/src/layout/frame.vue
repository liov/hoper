<template>
  <iframe :src="frameSrc" class="frame" frameborder="0" />
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();
const frameSrc = computed(() => (route.query?.url as string) ?? "");
</script>

<style scoped>
.frame {
  width: 100%;
  height: 100vh;
  border: none;
}
</style>
