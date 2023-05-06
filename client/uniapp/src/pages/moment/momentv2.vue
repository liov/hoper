<template>

        <uni-card @click="onClick"  :isFull="true" sub-title="副标题" extra="额外信息" :thumbnail="avatar" :cover="cover" >
            <!-- <image slot='cover' style="width: 100%;" :src="cover"></image> -->
            <text class="uni-body">这是一个带封面和操作栏的卡片示例，此示例展示了封面插槽和操作栏插槽的用法。</text>
            <view slot="actions" class="card-actions">
                <view class="card-actions-item" @click="actionsClick('分享')">
                    <uni-icons type="redo" size="18" color="#999"></uni-icons>
                    <text class="card-actions-item-text">分享</text>
                </view>
                <view class="card-actions-item" @click="actionsClick('点赞')">
                    <uni-icons type="heart" size="18" color="#999"></uni-icons>
                    <text class="card-actions-item-text">点赞</text>
                </view>
                <view class="card-actions-item" @click="actionsClick('评论')">
                    <uni-icons type="chatbubble" size="18" color="#999"></uni-icons>
                    <text class="card-actions-item-text">评论</text>
                </view>
            </view>
        </uni-card>

</template>

<script setup lang="ts">
const cover = 'https://web-assets.dcloud.net.cn/unidoc/zh/shuijiao.jpg';
const avatar = 'https://web-assets.dcloud.net.cn/unidoc/zh/unicloudlogo.png';

function onClick(e){
    console.log(e)
}

function actionsClick(text) {
    uni.showToast({
        title: text,
        icon: 'none'
    })
}
</script>


<style scoped>

</style>