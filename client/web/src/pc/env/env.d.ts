/// <reference types="vite/client" />
interface ImportMetaEnv {
  HOPRE_PLATFORM: string;
  HOPRE_STATIC_DIR: string;
  HOPRE_API_HOST: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
