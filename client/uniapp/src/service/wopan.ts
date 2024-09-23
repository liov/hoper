import { client } from 'diamond/wopan'
client.fetch = async function (url, method, headers, body) {
  return new Promise((resolve, reject) => {
    uni.request({
      url,
      method,
      header: headers,
      data: body,
      success: function (res) {
        resolve(res)
      },
      fail: function (err) {
        reject(err)
      },
    })
  })
}
client.setToken(uni.getStorageSync('accessToken'), uni.getStorageSync('accessToken'))
client.psToken = uni.getStorageSync('psToken')
client.proxy = 'http://localhost:8080'
console.log(client)
export default client
