import * as wopan from 'diamond/es/wopan'

export interface FileNode {
  parent: FileNode
  file: wopan.File
  subFiles: FileNode[]
  pageNo: number
  pageSize: number
  hasMore: boolean
}
