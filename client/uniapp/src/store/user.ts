import { ObjMap } from '@/utils/user'
import { defineStore } from 'pinia'

import type { UserBase } from '@gen/pb/user/user.model'
import { httpclient } from '@/api/common'
import { CommonResp } from '@hopeio/utils/types'
import UserService from '@/api/user'

export interface UserState {
  auth: any
  token: string
  userCache: Map<number, any>
}

const state: UserState = {
  auth: null,
  token: '',
  userCache: new Map<number, any>(),
}
const tokenKey = 'token'
const getters = {
  getUser: (state: UserState) => {
    return (id: number): UserBase => state.userCache.get(id)
  },
}

const actions = {
  async getAuth() {
    if (state.auth) return
    const token = uni.getStorageSync(tokenKey)
    if (token) {
      state.token = token
      const { data } = await httpclient.get<CommonResp<any>>(`/api/auth`)
      // 跟后端的初始化配合
      if (data.code === 0) state.auth = data.data
    }
  },
  async login(params) {
    try {
      const { code, data } = await UserService.login(params)

      if (code !== 0) {
        throw new Error('Bad credentials')
      }

      state.auth = data.user
      state.token = data.token
      uni.setStorageSync(tokenKey, data.token)
      httpclient.defaults.header.Authorization = data.token
      await uni.navigateTo({ url: '/' })
    } catch (error: any) {
      console.log(error)
    }
  },
  async signup(params) {
    try {
      const { data } = await UserService.signup(params)
      state.auth = data.user
      state.token = data.token
      uni.setStorageSync(tokenKey, data.token)
      httpclient.defaults.header.Authorization = data.token
      await uni.navigateTo({ url: '/' })
    } catch (error: any) {
      if (error.response && error.response.status === 401) {
        throw new Error('Bad credentials')
      }
      console.log(error)
    }
  },
  async appendUsersById(ids: number[]) {
    const noExistsId: number[] = []
    ids.forEach((value) => {
      if (!state.userCache.has(value)) noExistsId.push(value)
    })
    if (noExistsId.length > 0) {
      const data = await UserService.baseUserList(noExistsId)
      if (data.code && data.code !== 0)
        await uni.showToast({ title: data.msg, icon: 'error', duration: 1000 })
      else this.appendUsers(data.data.list)
    }
  },
  appendUsers(users) {
    for (const user of users) {
      state.userCache.set(user.id, user)
    }
  },
}

export const useUserStore = defineStore('user', {
  state: () => state,
  getters,
  actions,
})
