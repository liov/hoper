
console.log(process.env.VITE_STATIC_DIR);
console.log(process.env.VITE_API_HOST);
export const STATIC_DIR = process.env.VITE_STATIC_DIR;
export const API_HOST = process.env.VITE_API_HOST;


enum Environment {
  PROD = "prod",
  DEV = "dev",
}
