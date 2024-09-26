<route lang="json5">
{
style: {
navigationStyle: 'custom',
navigationBarTitleText: 'file-view',
},
}
</route>

<template>
 <view style="height: 100%;">
   <wd-navbar
     fixed
     left-text="返回"
     left-arrow
     @click-left="handleClickLeft"
     right-icon="tune"
     placeholder
   >
     <template v-slot:title>
       <text>文件列表</text>
     </template>
     <template v-slot:right>
       <wd-button @click="deleteFile" custom-class="title-right" type="primary" size="small" round>
         <wd-icon
           name="delete"
         />
       </wd-button>

     </template>
   </wd-navbar>
   <view @touchstart="touchStart" @touchend="touchEnd" style="height: 90%;">
   <image
     v-if="file.fileType == '1'"
     style="width: 100%; height: 90%;"
     :src="file.previewUrl"
     mode="aspectFit"
   ></image>
   <video id="myVideo"   v-if="file.fileType != '1'" style="width: 100%; height: 90%;" :src="file.previewUrl" controls></video>
   </view>
   <wd-tabbar v-model="tabbar">
     <wd-tabbar-item >
       <template #icon>
         <text>{{index+'/'+wopanStore.curDir.subFiles.length}}</text>
       </template>
     </wd-tabbar-item>
     <wd-tabbar-item>
       <template #icon>
         <wd-button @click="deleteFile" custom-class="title-right" type="primary" size="small" round>
           <wd-icon
             name="delete"
           />
         </wd-button>
       </template>
     </wd-tabbar-item>
   </wd-tabbar>
 </view>
</template>

<script setup lang="ts">
import { useWopanStore } from '@/store/wopan'
import {onLoad} from "@dcloudio/uni-app";
import {storeToRefs} from "pinia";
import * as wopan from "diamond/wopan";
const wopanStore = useWopanStore()
defineOptions({
  name: 'FileView',
})
const tabbar = ref(1)
let index = ref(0)
const  file = ref()
onLoad((options) => {
  index.value = parseInt(options.index)
  file.value =wopanStore.curDir.subFiles[index.value].file
  console.log("file-view",file.value)
})

function handleClickLeft() {
  uni.navigateBack()
}
const {proxy}= getCurrentInstance()
async function deleteFile(){
  console.log("deleteFile")
   await wopanStore.deleteCurDirFile(index.value);
  if(index.value > wopanStore.curDir.subFiles.length-1) {
    const dirId = wopanStore.curDir.file.id
    wopanStore.curDir = wopanStore.curDir.parent
    await wopan.DeleteFile(wopan.SpaceType.Private, [dirId], null)
    uni.redirectTo({url: '/pages/wopan/list'})
  }
  if(index.value == wopanStore.curDir.subFiles.length-1){
    index.value = index.value-1
  }
  file.value = wopanStore.curDir.subFiles[index.value].file
}
let currentX = 0
let currentY = 0
function touchStart(e:any){
  console.log(e)
  currentX = e.changedTouches[0].clientX
  currentY = e.changedTouches[0].clientY
}

function touchEnd(e:any){
  console.log(e)
 const moveX = currentX- e.changedTouches[0].clientX
  const moveY = currentY - e.changedTouches[0].clientY
  let add = false
  let sub = false
  if(moveX<0 || moveY<0){
    sub = moveX+moveY<0
  }
  if(moveX>0 || moveY>0){
    add = moveX+moveY>0
  }
  if(index.value<wopanStore.curDir.subFiles.length-1 && add){
    index.value = index.value+1
    file.value = wopanStore.curDir.subFiles[index.value].file
    return
  }
  if(index.value>0 && sub){
    index.value = index.value-1
    file.value = wopanStore.curDir.subFiles[index.value].file
  }
}
</script>

<style scoped lang="scss">
/* #ifdef H5 || APP-VUE */
::v-deep
  /* #endif */
uni-page-body,page {
  height: 100%;
}
.wot-theme-light{
  height: 100%;
}
</style>
