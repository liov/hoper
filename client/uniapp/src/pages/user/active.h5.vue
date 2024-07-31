<route lang="json5">
{
  style: {
    navigationStyle: 'custom',
    navigationBarTitleText: '激活',
  },
}
</route>

<template></template>

<script setup lang="ts">
import UserService from '@/service/user'
import { useUserStore } from '@/store/user'
import { useRouter, useRoute } from 'vue-router'

const store = useUserStore()
const router = useRouter()
const route = useRoute()

UserService.active(route.params.id, route.params.secret)
</script>

<style scoped></style>
