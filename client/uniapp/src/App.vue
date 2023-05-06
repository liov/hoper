<script setup lang="ts">
import {onLaunch, onShow, onHide, onLoad} from "@dcloudio/uni-app";
import {nextTick} from "vue";
import {globalStore} from "@/store";
onLaunch(() => {
  console.log("App Launch");
   globalStore.platform =  uni.getSystemInfoSync().platform
    console.log(globalStore.platform)
});

onLoad(()=>{
    getApp().globalData!.text = 'test text'
    const { text } = getApp().globalData;
})

onShow(() => {
  console.log("App Show");
});
onHide(() => {
  console.log("App Hide");
});

</script>
<style lang="scss">
@import '@dcloudio/uni-ui/lib/uni-scss/index.scss';
</style>
