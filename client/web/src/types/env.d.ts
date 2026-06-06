/// <reference types="vite/client" />

interface ImportMetaEnv extends ViteEnv{
  VITE_STATIC_DIR: string;
  VITE_API_HOST: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

declare const __APP_PLATFORM__: string
