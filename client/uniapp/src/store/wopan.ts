import wopanClient from '@/service/wopan'
import * as wopan from 'diamond/wopan'
import { defineStore } from 'pinia'

interface FileNode {
  parent: wopan.File
  file: wopan.File
  subFiles: FileNode[]
  pageNo: number
  pageSize: number
}

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
const rootFile = {
  parent: null,
  file: {
    fid: '',
    id: '',
    name: '',
    type: 1,
  },
  subFiles: [],
  pageNo: 0,
  pageSize: 20,
}
const state: WopanState = {
  accessToken: uni.getStorageSync(accessTokenKey),
  phone: uni.getStorageSync(phoneKey),
  psToken: uni.getStorageSync(psTokenKey),
  refreshToken: uni.getStorageSync(refreshTokenKey),
  file: rootFile,
  curDir: rootFile,
}

const getters = {}

const actions = {
  async PcWebLogin(params) {
    try {
      await wopan.PcWebLogin(params.phone, params.password)
    } catch (error: any) {
      console.log(error)
    }
  },
  async PcLoginVerifyCode(params) {
    try {
      const res = await wopan.PcLoginVerifyCode(params.phone, params.password, params.messageCode)
      uni.setStorageSync(accessTokenKey, res.access_token)
      uni.setStorageSync(refreshTokenKey, res.refresh_token)
      await uni.navigateTo({ url: '/wopan/list' })
    } catch (error: any) {
      console.log(error)
    }
  },
  async AppLoginByMobile(params) {
    console.log('params', params)
    try {
      console.log('params', params)
      const res = await wopan.AppLoginByMobile(params.phone, params.smsCode)
      state.accessToken = res.access_token
      state.refreshToken = res.refresh_token
      uni.setStorageSync(accessTokenKey, res.access_token)
      uni.setStorageSync(refreshTokenKey, res.refresh_token)
    } catch (error: any) {
      console.log(error)
    }
  },
  async VerifySetPwd() {
    try {
      const res = await wopan.VerifySetPwd()
      if (res.verifyResult === '01') {
        await uni.navigateTo({ url: '/pages/wopan/login?setpwd=1' })
      }
    } catch (e) {
      uni.navigateTo({ url: '/pages/wopan/login' })
      console.log(e)
    }
  },
  async FileList() {
    try {
      const res = await wopan.QueryAllFiles(
        wopan.SpaceType.Private,
        '',
        state.curDir.pageNo,
        state.curDir.pageSize,
        wopan.SortType.NameAsc,
        '',
      )
      state.curDir.subFiles.push(
        res.files.map(
          (file: wopan.File): FileNode => ({
            parent: state.curDir,
            file,
            subFiles: [],
            pageNo: 0,
            pageSize: 20,
          }),
        ),
      )
    } catch (error: any) {
      console.log(error)
    }
  },
}

export const useWopanStore = defineStore('wopan', {
  state: () => state,
  actions,
  getters,
})
