import { defineConfig } from "vite";
import babel from "vite-plugin-babel";

// https://vitejs.dev/config/
export default defineConfig({
  esbuild: {
    jsx: "automatic",
    jsxFactory: "el",
    jsxImportSource: "redom",
  },
  plugins: [
    babel({
      babelConfig: {
        babelrc: false,
        configFile: false,
        plugins: [
          "babel-plugin-transform-redom-jsx",
          [
            "transform-react-jsx",
            {
              pragma: "el",
            },
          ],
        ],
      },
    }),
  ],
});
