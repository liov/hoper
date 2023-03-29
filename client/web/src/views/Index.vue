<template>
  <div>
    <a-carousel autoplay arrows effect="fade">
      <template #prevArrow>
        <div class="custom-slick-arrow" style="left: 10px; z-index: 1">
          <left-circle-outlined />
        </div>
      </template>
      <template #nextArrow>
        <div class="custom-slick-arrow" style="right: 10px">
          <right-circle-outlined />
        </div>
      </template>
      <div><h3>建设中，我觉得你看瞬间比较好</h3></div>
      <div><h3>建设中，我觉得你去聊天比较好</h3></div>
      <div><h3>建设中，我觉得你写博客比较好</h3></div>
      <div><h3>建设中，我觉得你比较好看</h3></div>
    </a-carousel>

    <a-row>
      <a-col :span="6"> hoper </a-col>
      <a-col :span="12" style="text-align: center">
        <div>
          <iframe
            src="https://m.hoper.xyz"
            height="812"
            width="375"
            frameBorder="0"
          ></iframe>
        </div>
      </a-col>
      <a-col :span="6"> with </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { LeftCircleOutlined, RightCircleOutlined } from "@ant-design/icons-vue";
</script>
<style lang="less" scoped>
div .ant-carousel {
  :deep(div) {
    .slick-slide {
      text-align: center;
      height: 160px;
      line-height: 160px;
      background: #364d79;
      overflow: hidden;
    }
    .custom-slick-arrow {
      width: 25px;
      height: 25px;
      font-size: 25px;
      color: #fff;
      background-color: rgba(31, 45, 61, 0.11);
      opacity: 0.3;
    }
    .custom-slick-arrow:before {
      display: none;
    }
    .custom-slick-arrow:hover {
      opacity: 0.5;
    }
  }

  :deep(.slick-slide h3) {
    color: #fff;
  }
}
</style>
