export default defineNuxtConfig({
  ssr: true,
  modules: ['@nuxtjs/tailwindcss'],
  css: ['~/assets/css/tailwind.css'],
  runtimeConfig: {
    apiBaseUrl: process.env.NUXT_API_BASE_URL || process.env.NUXT_PUBLIC_API_BASE_URL || 'http://backend:8080',
    public: {
      apiBaseUrl: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8080',
      tmdbEnabled: process.env.NUXT_PUBLIC_TMDB_ENABLED === 'true'
    }
  }
})
