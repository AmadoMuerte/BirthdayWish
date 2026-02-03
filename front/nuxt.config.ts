// https://nuxt.com/docs/api/configuration/nuxt-config
import { fileURLToPath } from 'node:url'

export default defineNuxtConfig({
  devtools: { enabled: true },

  srcDir: 'src/',
  alias: {
    '@': fileURLToPath(new URL('./src', import.meta.url)),
  },

  dir: {
    pages:   'app/routes',
    layouts: 'app/layouts',
    app: 'app',
  },

  css: ['@/shared/styles/index.scss'],
})