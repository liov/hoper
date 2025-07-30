import {request} from '@hopeio/utils/uniapp'

export const Env: ImportMetaEnv = import.meta.env
console.log(Env)
export const STATIC_DIR = Env.VITE_STATIC_DIR
export const API_HOST = Env.VITE_API_HOST

const Prod = 'prod'

request.defaults.baseUrl = API_HOST
