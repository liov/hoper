import { readdirSync, statSync } from "node:fs";
import { join, relative, resolve } from "node:path";
import { defineConfig } from "vite";

function collectEntries(dir: string, root: string, entries: Record<string, string>): void {
  for (const name of readdirSync(dir)) {
    const fullPath = join(dir, name);
    const stat = statSync(fullPath);
    if (stat.isDirectory()) {
      collectEntries(fullPath, root, entries);
      continue;
    }
    if (!name.endsWith(".ts") || name.endsWith(".d.ts")) {
      continue;
    }
    const entryName = relative(root, fullPath).replace(/\.ts$/, "");
    entries[entryName] = fullPath;
  }
}

const srcRoot = resolve(__dirname, "src");
const input: Record<string, string> = {};
collectEntries(srcRoot, srcRoot, input);

export default defineConfig({
  build: {
    outDir: "dist",
    emptyOutDir: true,
    sourcemap: false,
    minify: false,
    rollupOptions: {
      input,
      external: ["@bufbuild/protobuf"],
      output: [
        {
          format: "es",
          entryFileNames: "[name].js",
          chunkFileNames: "chunk-[hash].js"
        },
        {
          format: "cjs",
          entryFileNames: "[name].cjs",
          chunkFileNames: "chunk-[hash].cjs",
          exports: "named"
        }
      ]
    }
  }
});
