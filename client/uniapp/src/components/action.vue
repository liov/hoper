<template>
    <view class="actions">
        <view class="actions-item" @click="actionsClick('分享')">
            <uni-icons type="redo" size="18" color="#999"></uni-icons>
            <text class="actions-item-text">分享</text>
        </view>

        <view class="actions-item" @click="actionsClick('评论')">
            <uni-icons type="chatbubble" size="18" color="#999"></uni-icons>
            <text class="actions-item-text">评论</text>
        </view>
        <view class="actions-item" @click="actionsClick('收藏')">
            <uni-icons type="star" size="18" color="#999"></uni-icons>
            <text class="actions-item-text">收藏</text>
        </view>
        <view class="actions-item" @click="actionsClick('点赞')">
            <uni-icons type="heart" size="18" color="#999"></uni-icons>
            <text class="actions-item-text">点赞</text>
        </view>
        <view class="actions-item" @click="actionsClick('···')">
            <uni-icons type="more-filled" size="18" color="#999"></uni-icons>
        </view>
    </view>
</template>

<script setup lang="ts">
function actionsClick(text){
    uni.showToast({
        title:text,
        icon:'none'
    })
}
</script>



<style scoped lang="scss">
.actions{
  display: flex;
  flex-direction: row;
  justify-content: space-around;
  align-items: center;
  height: 30px;
  border-top: 1px #eee solid;
  .actions-item{
    display: flex;
    flex-direction: row;
    align-items: center;
  }
  .actions-item-text{
    font-size: 12px;
    color: #666;
    margin-left: 5px;
  }
}
</style>