
console.log(process.env.HOPRE_STATIC_DIR);
console.log(process.env.HOPRE_API_HOST);
export const STATIC_DIR = process.env.HOPRE_STATIC_DIR;
export const API_HOST = process.env.HOPRE_API_HOST;


enum Platform {
  Weapp = "weapp",
  H5 = "h5",
}

enum Environment {
  PROD = "prod",
  DEV = "dev",
}
