<template>
  <div>
    <div style="marginbottom: 16px; text-align: center">
      <span style="marginright: 6px">Gutter (px): </span>
      <div style="width: 50%; margin: 0 auto">
        <a-slider
          v-model="gutterKey"
          :min="0"
          :max="Object.keys(gutters).length - 1"
          :marks="gutters"
          :step="null"
        />
      </div>
      <span style="marginright: 6px; text-align: center">列数:</span>
      <div style="width: 50%; margin: 0 auto">
        <a-slider
          v-model="colCountKey"
          :min="0"
          :max="Object.keys(colCounts).length - 1"
          :marks="colCounts"
          :step="null"
        />
      </div>
    </div>
    <a-row :gutter="gutters[gutterKey]">
      <a-col
        v-for="item in colCounts[colCountKey]"
        :key="item.toString()"
        :span="24 / colCounts[colCountKey]"
      >
        <div>
          <a-avatar size="large" :src="user.avatar_url" />
        </div>
        <div>{{ user.name }}</div>
        <div>{{ user.id }}</div>
        <div>{{ user.phone }}</div>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import axios from "axios";
import { useRoute } from "vue-router";

const gutterKey = 1;
const colCountKey = 2;

const gutters = ref({});
const arr = ref([8, 16, 24, 32, 40, 48]);
arr.value.forEach((value, i) => {
  gutters[i] = value;
});
const colCounts = ref({});
const arr1 = [2, 3, 4, 6, 8, 12];
arr1.forEach((value, i) => {
  colCounts[i] = value;
});

const route = useRoute();
const { data } = await axios.get(`/api` + route.path);

const user = ref(data.data);
</script>

<style scoped lang="less"></style>
