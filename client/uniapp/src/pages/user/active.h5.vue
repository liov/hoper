<template></template>

<script setup lang="ts">
import UserService from '@/api/user'
import { useUserStore } from '@/store/user'
import { useRouter, useRoute } from 'vue-router'
import i18n from '@/locale'

definePage({
  style: {
    navigationStyle: 'custom',
    navigationBarTitleText: i18n.global.t('page.active'),
  }
})

const store = useUserStore()
const router = useRouter()
const route = useRoute()

UserService.active(Number(route.params.id), String(route.params.secret))
</script>

<style scoped></style>
