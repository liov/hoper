import client from '@/service/wopan'
import * as wopan from 'diamond/wopan'
import { defineStore } from 'pinia'
import { FileNode } from '@/model/wopan'

export interface WopanState {
  file: FileNode
  curDir: FileNode
  accessToken: string
  refreshToken: string
  psToken: string
  phone: string
}

const accessTokenKey = 'accessToken'
const refreshTokenKey = 'refreshToken'
const psTokenKey = 'psToken'
const phoneKey = 'phone'
const rootFile: FileNode = {
  parent: null,
  file: {
    id: '0',
    name: 'root',
    type: 1,
  },
  subFiles: [],
  pageNo: 0,
  pageSize: 20,
  hasMore: true,
}
const state: WopanState = {
  accessToken: uni.getStorageSync(accessTokenKey),
  phone: uni.getStorageSync(phoneKey),
  psToken: uni.getStorageSync(psTokenKey),
  refreshToken: uni.getStorageSync(refreshTokenKey),
  file: rootFile,
  curDir: rootFile,
}
client.setToken(state.accessToken, state.refreshToken)
client.psToken = state.psToken
const getters = {
}

const actions = {
  async PcWebLogin(params) {
    await wopan.PcWebLogin(params.phone, params.password)
  },
  async PcLoginVerifyCode(params) {
    const res = await wopan.PcLoginVerifyCode(params.phone, params.password, params.messageCode)
    uni.setStorageSync(accessTokenKey, res.access_token)
    uni.setStorageSync(refreshTokenKey, res.refresh_token)
    await uni.navigateTo({ url: '/wopan/list' })
  },
  async AppLoginByMobile(params) {
    console.log('params', params)

    console.log('params', params)
    const res = await wopan.AppLoginByMobile(params.phone, params.smsCode)
    state.accessToken = res.access_token
    state.refreshToken = res.refresh_token
    uni.setStorageSync(accessTokenKey, res.access_token)
    uni.setStorageSync(refreshTokenKey, res.refresh_token)
  },
  async VerifySetPwd() {
    const res = await wopan.VerifySetPwd()
    if (res.verifyResult === '01') {
      await uni.navigateTo({ url: '/pages/wopan/login?setpwd=1' })
    }
  },
  async PrivateSpaceLogin(params) {
    const res = await wopan.PrivateSpaceLogin(params.passwd)
    state.psToken = res.psToken
    uni.setStorageSync(psTokenKey, res.psToken)
  },
  async FileList() {
    if (!state.curDir.hasMore) {
      return
    }
    const res = await wopan.QueryAllFiles(
      wopan.SpaceType.Private,
      state.curDir.file.id,
      state.curDir.pageNo,
      state.curDir.pageSize,
      wopan.SortType.NameAsc,
      '',
    )
    state.curDir.subFiles.push(
      ...res.files.map(
        (file: wopan.File): FileNode => {
          if(file.previewUrl===""){
            file.previewUrl = wopan.preview(file.fid)
          }
        return {
          parent: state.curDir,
          file,
          subFiles: [],
          pageNo: 0,
          pageSize: 50,
          hasMore: true,
        }
        },
      ),
    )
    if (res.files.length < state.curDir.pageSize) {
      state.curDir.hasMore = false
    }
    state.curDir.pageNo++
  },
  checkToken(){
    if (state.accessToken === '') {
      uni.navigateTo({
        url: '/pages/wopan/login',
      })
    }
    if (state.psToken === '') {
      uni.navigateTo({
        url: '/pages/wopan/login?psToken=1',
      })
    }
  },
  async deleteCurDirFile(index: number){
    if (state.curDir.file.id === '0') {
      return
    }
    const dirList= []
    const  fileList = []
    if (state.curDir.subFiles[index].file.type === 0) {
      dirList.push(state.curDir.subFiles[index].file.id)
    }else {
      fileList.push(state.curDir.subFiles[index].file.id)
    }
    console.log('deleteCurDirFile', dirList, fileList)
    await wopan.DeleteFile(wopan.SpaceType.Private, dirList, fileList)
    state.curDir.subFiles.splice(index, 1)
  }
}

export const useWopanStore = defineStore('wopan', {
  state: () => state,
  actions,
  getters,
})
