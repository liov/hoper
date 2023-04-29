<template>
    <uni-badge text="1"></uni-badge>
    <uni-badge text="2" type="success" ></uni-badge>
    <uni-badge text="3" type="primary" :inverted="true"></uni-badge>
    <uni-table border stripe emptyText="暂无更多数据" >
        <!-- 表头行 -->
        <uni-tr>
            <uni-th align="center">日期</uni-th>
            <uni-th align="center">姓名</uni-th>
            <uni-th align="left">地址</uni-th>
        </uni-tr>
        <!-- 表格数据行 -->
        <uni-tr>
            <uni-td>2020-10-20</uni-td>
            <uni-td>Jeson</uni-td>
            <uni-td>北京市海淀区</uni-td>
        </uni-tr>
        <uni-tr>
            <uni-td>2020-10-21</uni-td>
            <uni-td>HanMeiMei</uni-td>
            <uni-td>北京市海淀区</uni-td>
        </uni-tr>
        <uni-tr>
            <uni-td>2020-10-22</uni-td>
            <uni-td>LiLei</uni-td>
            <uni-td>北京市海淀区</uni-td>
        </uni-tr>
        <uni-tr>
            <uni-td>2020-10-23</uni-td>
            <uni-td>Danner</uni-td>
            <uni-td>北京市海淀区</uni-td>
        </uni-tr>

    </uni-table>
</template>

<script setup lang="ts">
//import {uniBadge} from '@dcloudio/uni-ui'
//import uniBadge from '@dcloudio/uni-ui/lib/uni-badge/uni-badge.vue'
</script>

<style>

</style>
