import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  compatibilityDate: '2026-05-01',
  devtools: { enabled: true },
  modules: ['@nuxtjs/robots', '@nuxtjs/sitemap'],

  components: [
    { path: '~/components', pathPrefix: false },
  ],

  css: ['~/assets/css/main.css'],

  vite: {
    plugins: [tailwindcss()],
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || '/api',
      siteUrl: process.env.NUXT_PUBLIC_SITE_URL || 'http://localhost:3000',
    },
  },

  app: {
    head: {
      htmlAttrs: { lang: 'en' },
      title: 'petal — a quiet mood journal',
      meta: [
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        {
          name: 'description',
          content:
            'petal is a calm, beautiful daily mood journal. Tap a petal to log how the day felt and watch your month bloom.',
        },
        { property: 'og:title', content: 'petal — a quiet mood journal' },
        { property: 'og:type', content: 'website' },
      ],
      link: [{ rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' }],
    },
  },

  routeRules: {
    '/': { prerender: true },
    '/pricing': { prerender: true },
    '/about': { prerender: true },
    '/login': { ssr: true },
    '/signup': { ssr: true },
    '/app/**': { ssr: false },
  },

  site: {
    url: process.env.NUXT_PUBLIC_SITE_URL || 'http://localhost:3000',
    name: 'petal',
  },

  nitro: {
    preset: 'node-server',
  },

  typescript: {
    strict: true,
  },
})
