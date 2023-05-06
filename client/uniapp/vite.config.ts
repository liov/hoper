import { defineConfig } from "vite";
import uni from "@dcloudio/vite-plugin-uni";
import * as path from 'path';
import inject from '@rollup/plugin-inject';


// https://vitejs.dev/config/
export default defineConfig({
  envDir: "./src/env",
  envPrefix: "HOPRE_",
  plugins: [
      uni()
  ],
  define: {
    'process.env.VUE_APP_TEST': JSON.stringify('test'),
  },
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
      },
    },
  },
});
