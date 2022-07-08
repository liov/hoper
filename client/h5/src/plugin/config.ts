const viteEnv: ImportMetaEnv = import.meta.env;
console.log(viteEnv);
export const STATIC_DIR = viteEnv.VITE_STATIC_DIR;
export const API_HOST = viteEnv.VITE_API_HOST;
