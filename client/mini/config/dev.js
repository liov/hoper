module.exports = {
  env: {
    NODE_ENV: '"development"',
    HOPRE_STATIC_DIR:'"https://static.hoper.xyz/hoper/"',
    HOPRE_API_HOST:'"https://api.hoper.xyz"'
  },
  defineConstants: {
    'process.env.HOPRE_STATIC_DIR': '"https://static.hoper.xyz/hoper/"',
    'process.env.HOPRE_API_HOST': '"https://api.hoper.xyz"'
  },
  mini: {},
  h5: {
/*    webpackChain (chain) {
      chain.plugin('sw')
        .use(require('workbox-webpack-plugin').GenerateSW, [{
          clientsClaim: true,
          skipWaiting: true,
        }])
    }*/
  }
}
