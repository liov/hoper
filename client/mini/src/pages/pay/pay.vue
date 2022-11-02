<template>
  <div class="pay">
    <nut-tabs
      v-model="active"
      @change="onChange"
    >
      <template #nav-left>
        <div style="width: 30%"></div>
      </template>
      <template #nav-right>
        <div style="width: 30%">
          <nut-icon name="add-o" color="#1989fa" size="1.5rem" />
        </div>
      </template>
      <nut-tabpane title="关注"> 内容1 </nut-tabpane>
      <nut-tabpane title="推荐"> 内容2 </nut-tabpane>
      <nut-tabpane title="最新"><MomentList v-if="active==2"/> </nut-tabpane>
    </nut-tabs>

    <nut-row type="flex" wrap="nowrap" class="row">
      <nut-col span="6"> span: 6 </nut-col>
      <nut-col span="12">span: 12 </nut-col>
      <nut-col span="6">span: 6</nut-col>
    </nut-row>
  </div>
</template>

<script setup lang="ts">
import { Toast } from "@/plugins/compatible/toast";
import MomentList from "@/pages/moment/list.vue";
import { defineAsyncComponent, ref } from "vue";
const active = ref(1);

function onChange(tab) {
  console.log(tab)
}
</script>

<style scoped lang="scss">
i {
  margin: 0.7rem;
}
.nut-tabpane{
  padding: 0
}
</style>
