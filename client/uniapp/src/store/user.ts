import { ObjMap } from '@/utils/user'
import { defineStore } from 'pinia'
import { API_HOST } from '@/env/config'
import type { User } from '@/model/user'
import request from '@/utils/request'

export interface UserState {
  info: any
  token: string
  userCache: Map<number, any>
}

export const state: UserState = {
  info: null,
  token: '',
  userCache: new Map<number, any>(),
}

const getters = {
  getUser: (state) => {
    return (id): User => state.userCache.get(id)
  },
}

const actions = {
  async getAuth() {
    if (state.info) return
    const token = uni.getStorageSync('token')
    if (token) {
      state.token = token
      const { data } = await request.get(API_HOST + `/api/v1/auth`)
      // 跟后端的初始化配合
      if (data.code === 0) state.info = data.data
    }
  },
  async login(params) {
    try {
      const { code, data } = await request.post('/api/v1/user/login', params)

      if (code !== 0) {
        throw new Error('Bad credentials')
      }

      state.info = data.user
      state.token = data.token
      uni.setStorageSync('token', data.token)
      request.defaults.header.Authorization = data.token
      await uni.navigateTo({ url: '/' })
    } catch (error: any) {
      console.log(error)
    }
  },
  async signup(params) {
    try {
      const {
        data: { details },
      } = await request.post('/api/v2/user', params, {
        header: {
          'content-type': 'application/json',
        },
      })
      state.info = details.user
      state.token = details.token
      uni.setStorageSync('token', details.token)
      request.defaults.header.Authorization = details.token
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
      const data = await request.post(`/api/v1/user/baseUserList`, {
        ids: noExistsId,
      })
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

export const useUserStore = defineStore({
  id: 'user',
  state: () => state,
  getters,
  actions,
})
