/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./template/**/*.templ",
    "./static/js/**/*.js"
  ],
  theme: {
    extend: {
      fontFamily: {
        primary: ['Mukta', 'sans-serif'],
        accent: ['Caveat', 'cursive'],
      },
    },
  },
  plugins: [
    require('daisyui'),
  ],
  daisyui: {
    themes: [
      {
        cmyk: {
          ...require("daisyui/src/theming/themes")["cmyk"],
          primary: "#E2344D",
          secondary: "#1783B5",
          success: "green",
        },
      },
    ],
    logs: true,
  },
};
