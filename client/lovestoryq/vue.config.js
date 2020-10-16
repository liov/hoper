module.exports = {
  devServer: {
    // overlay: { // 让浏览器 overlay 同时显示警告和错误
    //   warnings: true,
    //   errors: true
    // },
    // open: false, // 是否打开浏览器
    host: "hoper.xyz",
    port: "80", // 代理断就
    //https: true,
    // hotOnly: false, // 热更新
    proxy: {
      "/api": {
        target: "https://106.54.79.41", // 目标代理接口地址
        // ws: true, // 是否启用websockets
        changeOrigin: true // 开启代理，在本地创建一个虚拟服务端
        // pathRewrite: {"^/api": "/"}
      },
      "/static": {
        target: "https://106.54.79.41", // 目标代理接口地址
        // ws: true, // 是否启用websockets
        changeOrigin: true // 开启代理，在本地创建一个虚拟服务端
        // pathRewrite: {"^/api": "/"}
      }
    }
  }
};
