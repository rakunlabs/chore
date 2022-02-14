import * as path from "path";
import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import preprocess from "svelte-preprocess";

// https://vitejs.dev/config/
export default defineConfig({
  base: "./",
  plugins: [svelte({
    preprocess: preprocess({ name: "scss" }),
    compilerOptions: {
      dev: process.env.NODE_ENV == "production" ? false : true,
    },
  })],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
    },
  },
  build: {
    sourcemap: process.env.NODE_ENV == "production" ? (process.env.SOURCEMAP == "true" ? true : false) : true,
  },
  mode: process.env.NODE_ENV == "production" ? "" : "development",
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        secure: true,
        ws: true,
        followRedirects: false,
      },
    },
    port: process.env.PORT ?? 3000,
  },
});
