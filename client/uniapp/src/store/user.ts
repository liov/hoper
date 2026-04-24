import { ObjMap } from '@/utils/user'
import { defineStore } from 'pinia'

import type { UserBase } from '@gen/pb/user/user.model'
import { httpclient } from '@/api/common'
import { CommonResp } from '@hopeio/utils/types'
import UserService from '@/api/user'
import { LoginReq, SignupReq } from '@gen/pb/user/user.service'

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
export const tokenKey = 'token'
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
      const data = await UserService.auth()
      // 跟后端的初始化配合
       state.auth = data
    }
  },
  async login(params: LoginReq) {
    try {
      const data = await UserService.login(params)
      state.auth = data.user
      state.token = data.token
      uni.setStorageSync(tokenKey, data.token)
      httpclient.defaults.header.Authorization = data.token
      await uni.switchTab({ url: '/pages/moment/moment_list' })
    } catch (error) {
      console.log(error)
    }
  },
  async signup(params: SignupReq) {
    const data = await UserService.signup(params)
    const tip = typeof data === 'string' ? data : '注册成功'
    await uni.showToast({ title: tip, icon: 'success', duration: 2500 })
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
  logout() {
    state.auth = null
    state.token = ''
    uni.removeStorageSync(tokenKey)
    httpclient.defaults.header.Authorization = ''
    uni.reLaunch({ url: '/pages/user/login' })
  },
}

export const useUserStore = defineStore('user', {
  state: () => state,
  getters,
  actions,
})
