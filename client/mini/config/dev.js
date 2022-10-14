module.exports = {
  env: {
    NODE_ENV: '"development"',
    VITE_STATIC_DIR:'"https://static.hoper.xyz/hoper/"',
    VITE_API_HOST:'"https://api.hoper.xyz"'
  },
  defineConstants: {
    'process.env.VITE_STATIC_DIR': '"https://static.hoper.xyz/hoper/"',
    'process.env.VITE_API_HOST': '"https://api.hoper.xyz"'
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
