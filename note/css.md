# display: flex 无法居中

justify-content: center;

# 深度作用选择器
/deep/ 改为 ::v-deep


你很可能会遇到的问题

vue组件编译后，会将 template 中的每个元素加入 [data-v-xxxx] 属性来确保 style scoped 仅本组件的元素而不会污染全局，但如果你引用了第三方组件：



默认只会对组件的最外层（div）加入这个 [data-v-xxxx] 属性，但第二层开始就没有效果了。如图所示： 第一层还有 data-v-17bb9a05, 但第二层的 .weui-cells 就没有了。



（这是 <group />组件的源码： https://github.com/airyland/vux/blob/v2/src/components/group/index.vue ）



所以，如果你期待通过如下方式修改 weui-cells 的样式。是没有任何效果的：

<style scoped>
    .fuck .weui-cells {
        // ...
    }
</style>


这是因为，所有的scoped中的css最终编译出来都会变成这样：

.fuck[data-v-17bb9a05] .weui-cells[data-v-17bb9a05]


解决方法一：除非你将 scoped 移除，或者新建一个没有 scoped 的 style（一个.vue文件允许多个style）

<style scoped>
    .fuck {
        // ...
    }
</style>

<style>
    .fuck .weui-cells {
        // ...
    }
</style>


解决方法二：深度作用选择器 >>>

（注意，只作用于css）

.fuck >>> .weui-cells {
// ...
}
但如果是sass/less的话可能无法识别，这时候需要使用 /deep/ 选择器。

<style lang="scss" scoped>
.select {
  width: 100px;

  /deep/ .el-input__inner {
    border: 0;
    color: #000;
  }
}
</style>

the >>> and /deep/ combinators have been deprecated. Use :deep() instead.