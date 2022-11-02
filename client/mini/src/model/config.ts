
console.log(process.env.HOPRE_STATIC_DIR);
console.log(process.env.HOPRE_API_HOST);
export const STATIC_DIR = process.env.HOPRE_STATIC_DIR;
export const API_HOST = process.env.HOPRE_API_HOST;


export const enum Platform {
  Weapp = "weapp",
  H5 = "h5",
}

export const enum Environment {
  PROD = "prod",
  DEV = "dev",
}
