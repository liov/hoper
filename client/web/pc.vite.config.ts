import { fileURLToPath, URL } from "url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import vueJsx from "@vitejs/plugin-vue-jsx";
import Components from "unplugin-vue-components/vite";
import { AntDesignVueResolver } from "unplugin-vue-components/resolvers";
import { VitePWA } from "vite-plugin-pwa";
import Unocss from "unocss/vite";
import { presetUno, presetAttributify, presetIcons } from "unocss";

// TS1259: Module '"path"' can only be default-imported using the 'esModuleInterop' flag
import * as path from "path";
import fs from "fs";
import wasm from "vite-plugin-wasm";
import { viteCommonjs } from "@originjs/vite-plugin-commonjs";

const lessVar = path.resolve(__dirname, "src/pc/assets/var.less");

fs.copyFileSync("src/pc/index.html", "index.html");

// https://vitejs.dev/config/
export default defineConfig({
  envDir: "./src/pc/env",
  envPrefix: "HOPRE_",
  server: {
    port: 80,
    strictPort: true, // 严格端口 true:如果端口已被使用，则直接退出，而不会再进行后续端口的尝试。
    /**
     * @description 解决chrome设置origin:*也跨域机制,代理/api前缀到服务基地址
     * 最终的地址会将axios设置的baseUrl:/api代理拼接成[target][/api],然后通过rewrite重写掉/api为'' 这样就是我们真实的基地址了
     */
    proxy: {
      "/static": {
        target: "https://static.hoper.xyz", // 接口基地址
        rewrite: (path) => {
          console.log(path); // 打印[/api/userInfo] 这就是http-proxy要请求的url,我们基地址实际是没有/api 所以replace掉
          return path.replace(/^\/static/, "/hoper");
        },
      },
    },
  },
  optimizeDeps: {
    include: [
      "qs",
      "mitt",
      "xlsx",
      "dayjs",
      "axios",
      "pinia",
      "echarts",
      "esm-dep > cjs-dep",
    ],
  },
  plugins: [
    vue(),
    vueJsx(),
    viteCommonjs(),
    Components({
      resolvers: [
        AntDesignVueResolver({
          importStyle: false, // css in js
        }),
      ],
    }),
    ...VitePWA({ registerType: "autoUpdate", outDir: "dist/pc" }),
    //wasm(),
    //ViteRsw(),
    Unocss({
      // 使用Unocss
      presets: [presetUno(), presetAttributify(), presetIcons()],
      rules: [
        // 在这个可以增加预设规则, 也可以使用正则表达式
        [
          "p-c", // 使用时只需要写 p-c 即可应用该组样式
          {
            position: "absolute",
            top: "50%",
            left: "50%",
            transform: `translate(-50%, -50%)`,
          },
        ],
        [/^m-(\d+)$/, ([, d]) => ({ margin: `${parseInt(d) / 4}rem` })],
      ],
    }),
  ],
  worker: {
    format: "es",
    //plugins: [wasm()],
  },
  resolve: {
    alias: {
      // "@": fileURLToPath(new URL("./src", import.meta.url)),
      "@pc": fileURLToPath(new URL("./src/pc", import.meta.url)),
      "@generated": fileURLToPath(new URL("./generated", import.meta.url)),
      "@types": fileURLToPath(new URL("./types", import.meta.url)),
    },
  },
  build: {
    rollupOptions: {
      // https://rollupjs.org/guide/en/#outputmanualchunks
      output: {
        dir: "dist/pc",
        manualChunks: {},
      },
      plugins: [],
    },
    minify: "terser",
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true,
      },
    },
    dynamicImportVarsOptions: {
      include: ["./src/pc/**/*.ts"],
    },
  },
  css: {
    preprocessorOptions: {
      less: {
        javascriptEnabled: true,
        additionalData: `@import "${lessVar}";`,
      },
    },
  },
});
