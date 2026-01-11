import { heroui } from '@heroui/react'

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{js,ts,jsx,tsx}',
    // HeroUI theme 可能在不同位置，使用通配符匹配
    './node_modules/@heroui/theme/dist/**/*.{js,ts,jsx,tsx}',
    './node_modules/@heroui/*/node_modules/@heroui/theme/dist/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {},
  },
  darkMode: 'class',
  plugins: [
    heroui({
      themes: {
        light: {
          colors: {
            background: '#FAFAF8',
            foreground: '#2D2D2A',
            content1: '#F5F5F3',
            content2: '#EFEFE9',
            content3: '#E8E8E2',
            content4: '#E0E0DA',
            default: {
              50: '#FAFAF8',
              100: '#F5F5F3',
              200: '#EBEBEA',
              300: '#DCDCDA',
              400: '#A8A8A5',
              500: '#737370',
              600: '#5C5C58',
              700: '#454542',
              800: '#2D2D2A',
              900: '#1C1C1A',
              foreground: '#2D2D2A',
              DEFAULT: '#737370',
            },
          },
        },
        dark: {
          colors: {
            background: '#1C1C1A',
            foreground: '#EDEDE8',
            content1: '#252523',
            content2: '#2D2D2A',
            content3: '#353532',
            content4: '#3D3D3A',
            default: {
              50: '#252523',
              100: '#2D2D2A',
              200: '#3D3D3A',
              300: '#525250',
              400: '#6B6B68',
              500: '#8A8A86',
              600: '#A8A8A5',
              700: '#C8C8C4',
              800: '#E0E0DC',
              900: '#EDEDE8',
              foreground: '#EDEDE8',
              DEFAULT: '#8A8A86',
            },
          },
        },
      },
    }),
  ],
}
