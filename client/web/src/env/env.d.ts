/// <reference types="vite/client" />
interface ImportMetaEnv {
  HOPRE_STATIC_DIR: string;
  HOPRE_API_HOST: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

declare const __APP_PLATFORM__: string
