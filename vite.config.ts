import { defineConfig } from "vite";
import tailwindcss from "tailwindcss";

export default defineConfig({
  build: {
    lib: {
      // Could also be a dictionary or array of multiple entry points
      entry: ["client-components/index.ts"],
      name: "bundle",
      // the proper extensions will be added
      fileName: "bundle",
    },
  },
  css: {
    postcss: {
      plugins: [tailwindcss],
    },
  },
});
