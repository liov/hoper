import wopanClient from '@/service/wopan'
import * as wopan from 'diamond/wopan'
import { defineStore } from 'pinia'
interface File {
  fid: string
  name: string
  id: string
  type: number
  subFiles: Files
}

interface Files {
  parent: File
  files: File[]
  pageNo: number
  pageSize: number
}

export interface WopanState {
  files: Files
  accessToken: string
  refreshToken: string
  psToken: string
  phone: string
}

export const state: WopanState = {
  files: null,
}

const getters = {}

const actions = {
  async PcWebLogin(params) {
    try {
      const res = await wopan.PcWebLogin(params.phone, params.password)
    } catch (error: any) {
      console.log(error)
    }
  },
  async PcLoginVerifyCode(params) {
    try {
      const res = await wopan.PcLoginVerifyCode(params.phone, params.password, params.messageCode)
      uni.setStorageSync('accessToken', res.access_token)
      uni.setStorageSync('refreshToken', res.refresh_token)
      await uni.navigateTo({ url: '/wopan/list' })
    } catch (error: any) {
      console.log(error)
    }
  },
}

export const useWopanStore = defineStore({
  id: 'wopan',
  state: () => state,
  getters,
  actions,
})
