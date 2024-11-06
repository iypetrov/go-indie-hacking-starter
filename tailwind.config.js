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
      },
    },
  },
  plugins: [
    require('daisyui'),
  ],
  daisyui: {
    themes: ["dark"],
    logs: true,
  },
};
