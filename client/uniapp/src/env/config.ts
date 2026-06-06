
export const Env: ImportMetaEnv = import.meta.env
console.log(Env)
export const STATIC_DIR = Env.VITE_STATIC_DIR

const Prod = 'prod'
