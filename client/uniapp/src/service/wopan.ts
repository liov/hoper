import { client } from 'diamond/wopan'
client.fetch = async function (url, method, headers, body) {
  return new Promise((resolve, reject) => {
    uni.request({
      url,
      method,
      header: headers,
      data: body,
      timeout: 20000,
      success: function (res) {
        resolve(res)
      },
      fail: function (err) {
        reject(err)
      },
    })
  })
}
client.failCallback = function (msg:string) {
  if (msg.startsWith('request failed with rsp_code: 1001')) {
    uni.navigateTo({ url: '/pages/wopan/login' })
  } else if (msg.startsWith('request failed with rsp_code: 110001')) {
    uni.navigateTo({ url: '/pages/wopan/login?psToken=1' })
  }
  uni.showToast({
    icon:'none',
    title: msg,
  })
  throw new Error(msg)
}
client.proxy = 'http://localhost:8080'
console.log(client)
export default client
