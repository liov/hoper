使用 scoped 后，父组件的样式将不会渗透到子组件中。

例如（无效）：
```jsx
<template>
  <div id="app">
    <el-input  class="text-box" v-model="text"></el-input>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      text: 'hello'
    };
  }
};
</script>

<style lang="less" scoped>
.text-box {
   input {
    width: 166px;
    text-align: center;
  }
}
</style>
```
解决方法:

使用深度作用选择器 /deep/
```jsx
<template>
  <div id="app">
    <el-input v-model="text" class="text-box"></el-input>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      text: 'hello'
    };
  }
};
</script>

<style lang="less" scoped>
.text-box {
  /deep/ input {
    width: 166px;
    text-align: center;
  }
}
</style>
```