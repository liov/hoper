export const Env: ImportMetaEnv = import.meta.env;
console.log(Env);
export const STATIC_DIR = Env.VITE_STATIC_DIR;
export const API_HOST = Env.VITE_API_HOST;
export const APP_PLATFORM = __APP_PLATFORM__;
