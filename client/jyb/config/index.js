const config = {
  projectName: 'jyb',
  date: '2020-11-11',
  designWidth: 750,
  deviceRatio: {
    640: 2.34 / 2,
    750: 1,
    828: 1.81 / 2
  },
  sourceRoot: 'src',
  outputRoot: 'dist',
  plugins: [],
  defineConstants: {
  },
  copy: {
    patterns: [
    ],
    options: {
    }
  },
  framework: 'vue3',
  mini: {
    postcss: {
      pxtransform: {
        enable: true,
        config: {

        }
      },
      url: {
        enable: true,
        config: {
          limit: 1024 // 设定转换尺寸上限
        }
      },
      cssModules: {
        enable: false, // 默认为 false，如需使用 css modules 功能，则设为 true
        config: {
          namingPattern: 'module', // 转换模式，取值为 global/module
          generateScopedName: '[name]__[local]___[hash:base64:5]'
        }
      }
    }
  },
  h5: {
    publicPath: '/',
    staticDirectory: 'static',
    postcss: {
      autoprefixer: {
        enable: true,
        config: {
        }
      },
      cssModules: {
        enable: false, // 默认为 false，如需使用 css modules 功能，则设为 true
        config: {
          namingPattern: 'module', // 转换模式，取值为 global/module
          generateScopedName: '[name]__[local]___[hash:base64:5]'
        }
      },
    },
    esnextModules: ['taro-ui-vue3'],
    webpackChain(chain) {
      chain.resolve.alias
        .set(
          '@tarojs/components$',
          '@tarojs/components/dist-h5/vue3/index.js'
        )
    },
    devServer: {
      // overlay: { // 让浏览器 overlay 同时显示警告和错误
      //   warnings: true,
      //   errors: true
      // },
      // open: false, // 是否打开浏览器
      //host: "liov.xyz",
      //port: "80", // 代理断就
      //https: true,
      // hotOnly: false, // 热更新
      proxy: {
        "/api": {
          target: "https://hoper.xyz", // 目标代理接口地址
          // ws: true, // 是否启用websockets
          changeOrigin: true // 开启代理，在本地创建一个虚拟服务端
          // pathRewrite: {"^/api": "/"}
        },
        "/static": {
          target: "https://hoper.xyz", // 目标代理接口地址
          // ws: true, // 是否启用websockets
          changeOrigin: true // 开启代理，在本地创建一个虚拟服务端
          // pathRewrite: {"^/api": "/"}
        }
      }
    }
  }
}

module.exports = function (merge) {
  if (process.env.NODE_ENV === 'development') {
    return merge({}, config, require('./dev'))
  }
  return merge({}, config, require('./prod'))
}
