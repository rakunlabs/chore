/* eslint-disable @typescript-eslint/no-var-requires */
const tailwindcssTypography = require("@tailwindcss/typography");
const tailwindcssNesting = require("tailwindcss/nesting");
const tailwindcss = require("tailwindcss");
const autoprefixer = require("autoprefixer");

module.exports = {
  mode: "jit",
  content: ["./src/app.html", "./src/**/*.{svelte,js,ts,jsx,tsx}"],
  darkMode: "class", // or 'media' or 'class'
  theme: {
    extend: {},
  },
  variants: {
    extend: {},
  },
  plugins: [
    tailwindcssTypography,
    tailwindcssNesting,
    tailwindcss,
    autoprefixer,
  ],
};
