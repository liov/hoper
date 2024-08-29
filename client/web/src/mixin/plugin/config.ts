export const Env: ImportMetaEnv = import.meta.env;
console.log(Env);
export const STATIC_DIR = Env.HOPRE_STATIC_DIR;
export const API_HOST = Env.HOPRE_API_HOST;
export const APP_PLATFORM = __APP_PLATFORM__;
