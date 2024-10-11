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
   <view @touchstart="touchStart" @touchend="touchEnd" style="height: 80%;">
   <image
     v-if="file.fileType == '1'"
     style="width: 100%; height: 100%;"
     :src="file.previewUrl"
     mode="aspectFit"
   ></image>
   <video id="myVideo"   v-if="file.fileType != '1'" style="width: 100%; height: 90%;" :src="file.previewUrl" controls></video>
   </view>
   <view style="text-align: center"><text>{{fileName}}</text></view>
   <wd-tabbar v-model="tabbar" fixed>
     <wd-tabbar-item >
       <template #icon>
         <text>{{index+1+'/'+ total}}</text>
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
import * as wopan from "diamond/wopan";
import { useToast } from 'wot-design-uni'
import {FileNode} from "@/model/wopan";
const wopanStore = useWopanStore()
defineOptions({
  name: 'FileView',
})
const toast = useToast()
const tabbar = ref(1)
let index = ref(0)
const  file: Ref<wopan.File> = ref()
const fileName = computed(()=>{
  const len = 40
  const n = Math.floor(file.value.name.length/len)
  let fileName = file.value.name.slice(0,len)
  for(let i=1;i<n;i++){
    fileName+='\n'+file.value.name.slice(i*len,(i+1)*len)
  }
  if(len*n < file.value.name.length){
    fileName+='\n'+file.value.name.slice(n*len)
  }
  return fileName
})


const total = ref(wopanStore.curDir.subFiles.length)
onLoad((options) => {
  index.value = parseInt(options.index)
  file.value =wopanStore.curDir.subFiles[index.value].file
  console.log("file-view",file.value)
})

function handleClickLeft() {
  uni.navigateBack()
}
const {proxy}= getCurrentInstance()
let hasMore = true

async function deleteFile(){
  console.log("deleteFile")
   await wopanStore.deleteCurDirFile(index.value);
  total.value = wopanStore.curDir.subFiles.length
  if(index.value == wopanStore.curDir.subFiles.length && !hasMore){
    await wopanStore.FileList()
  }
  if(index.value == wopanStore.curDir.subFiles.length){
    if(hasMore){
      hasMore = false
      toast.show('没有更多')
    }
    if (index.value > 0 ){
      index.value = index.value-1
    }
  }
  file.value = wopanStore.curDir.subFiles[index.value].file
}
let currentX = 0
let currentY = 0
function touchStart(e:any){
  currentX = e.changedTouches[0].clientX
  currentY = e.changedTouches[0].clientY
}

async function touchEnd(e:any){
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
  if(index.value<wopanStore.curDir.subFiles.length && add){
    index.value = index.value+1
    if (index.value == wopanStore.curDir.subFiles.length && !hasMore){
      toast.show('没有更多')
      index.value = index.value-1
      return
    }
    await wopanStore.FileList()
    if(index.value == wopanStore.curDir.subFiles.length) {
      if (hasMore) {
        hasMore = false
        toast.show('没有更多')
        index.value = index.value-1
        return
      }
    }
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
