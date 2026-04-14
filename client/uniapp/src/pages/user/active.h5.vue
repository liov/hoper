<template></template>

<script setup lang="ts">
import UserService from '@/api/user'
import { useUserStore } from '@/store/user'
import { useRouter, useRoute } from 'vue-router'

definePage({
  style: {
    navigationStyle: 'custom',
    navigationBarTitleText: '激活',
  }
})

const store = useUserStore()
const router = useRouter()
const route = useRoute()

UserService.active(Number(route.params.id), String(route.params.secret))
</script>

<style scoped></style>
