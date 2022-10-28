
export function get(key: string): string | null {
  return wx.getStorageSync(key);
}
