import { defineConfig } from "vite";
import path from "path";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  define: {
    "process.env": process.env,
  },
  server: {
    port: 3000,
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@tasker/openapi": path.resolve(
        __dirname,
        "../../packages/openapi/src"
      ),
      "@tasker/zod": path.resolve(__dirname, "../../packages/zod/src"),
    },
  },
});
