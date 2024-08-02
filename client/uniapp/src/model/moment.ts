import type { ContentExt } from '@/model/content'

export interface Moment {
  id: number
  content: string
  images: string
  imagesUrls: string[]
  type: number
  userId: number
  user: User
  ext: ContentExt
  permission: number
  createdAt: string
}

export interface User {
  id: number
  name: string
  gender: number
  avatarUrl: string
}

export interface MomentList {
  total: number
  list: Moment[]
  users: User[]
}
