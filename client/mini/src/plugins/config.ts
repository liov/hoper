export const Env: NodeJS.ProcessEnv = process.env;
console.log(Env);
export const STATIC_DIR = Env.VITE_STATIC_DIR;
export const API_HOST = Env.VITE_API_HOST;


enum Environment {
  PROD = "prod",
  DEV = "dev",
}
