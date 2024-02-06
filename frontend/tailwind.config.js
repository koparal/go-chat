/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx}',
    './components/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    colors: {
      'dark-primary': '#131a1c',
      'dark-secondary': '#1b2224',
      red: '#e74c4c',
      green: '#6bb05d',
      blue: '#112750',
      grey: '#86898d',
      white: '#fff',
    },
  },
  plugins: [],
}
