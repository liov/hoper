import * as wopan from 'diamond/es/wopan'

export interface FileNode {
  parent: FileNode
  file: wopan.File
  subFiles: FileNode[]
  pageNo: number
  hasMore: boolean
  read: boolean
  deleted: boolean
}
