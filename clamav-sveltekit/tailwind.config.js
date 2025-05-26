/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        blue: {
          600: '#4361ee',
          700: '#3a0ca3'
        },
        cyan: {
          500: '#4cc9f0'
        },
        pink: {
          600: '#f72585'
        }
      },
      transitionProperty: {
        'height': 'height',
        'spacing': 'margin, padding'
      },
      rotate: {
        '30': '30deg'
      }
    }
  },
  plugins: []
}; 