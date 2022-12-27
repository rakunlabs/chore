import * as path from "path";
import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import preprocess from "svelte-preprocess";
import basicSsl from '@vitejs/plugin-basic-ssl';

// https://vitejs.dev/config/
export default defineConfig({
  base: process.env.NODE_ENV == "production" ? "./" : process.env.BASEPATH ?? "./",
  plugins: [
    svelte({
      preprocess: preprocess({ name: "scss" }),
      compilerOptions: {
        dev: process.env.NODE_ENV == "production" ? false : true,
      },
    }),
    !!process.env.HTTPS ? basicSsl() : null,
  ],
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
      [process.env.BASEPATH ? path.join(process.env.BASEPATH + "api") : "/api"]: {
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
