<route lang="json5">
{
  style: {
    navigationBarTitleText: '登录',
  },
}
</route>

<template>
  <div id="cf-turnstile"></div>
</template>

<script setup lang="ts">
import { dynamicLoadJs } from '@/utils/h5'

const turnstile = window.turnstile
const onloadTurnstileCallback = function () {
  turnstile.ready(function () {
    turnstile.render('#cf-turnstile', {
      sitekey: '0x4AAAAAAAgC9s4WZMlljGRg',
      callback: function (token) {
        console.log(`Challenge Success ${token}`)
      },
    })
  })

  if (!turnstile) {
    dynamicLoadJs(
      '//challenges.cloudflare.com/turnstile/v0/api.js?render=explicit"',
      onloadTurnstileCallback,
    )
  } else {
    onloadTurnstileCallback()
  }
}
</script>

<style scoped></style>
