<template>
  <view class="index">
    <NumberDisplay/>
    <NumberSubmit/>
    <AtButton type='primary' @click="handleClick">按钮文案</AtButton>
  </view>
</template>

<script>
import NumberDisplay from '../../components/NumberDisplay.vue'
import NumberSubmit from '../../components/NumberSubmit.vue'
import { AtButton } from 'taro-ui-vue3'

export default {
  name: 'Index',
  components: {
    NumberDisplay,
    NumberSubmit,
    AtButton
  }
}
</script>

<style>
.index {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
