// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  runtimeConfig: {
    public: {
      apiBase: "http://localhost:3000"
    }
  },
  vite: {
    server: {
      proxy: {
        '/rpc': 'http://localhost:8080',
        '/v1': {
          target: 'http://localhost:8080',
          cookieDomainRewrite: 'localhost'
        },
      }
    }
  },
  modules: [
    '@nuxtjs/color-mode',
    '@nuxtjs/tailwindcss',
    'nuxt-icon',
    'nuxt-headlessui',
    '@pinia/nuxt',
    '@vueuse/nuxt'
  ],
  colorMode: {
    classPrefix: 'theme-',
    classSuffix: '',
  },
  app: {
    head: {
      title: "IPCManView",
      htmlAttrs: {
        lang: "en-US",
      },
    },
  },
})
