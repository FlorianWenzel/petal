// Forwards every /api/* request to the backend API container/service.
// Runtime env so the same build works in dev, local docker, and prod.
//   NUXT_API_PROXY_TARGET=http://api:4000      (docker compose, Coolify)
//   NUXT_API_PROXY_TARGET=http://localhost:4000 (npm run dev:web)
export default defineEventHandler(async (event) => {
  const url = getRequestURL(event)
  if (!url.pathname.startsWith('/api/') && url.pathname !== '/api') return

  const target = process.env.NUXT_API_PROXY_TARGET || 'http://api:4000'
  const stripped = url.pathname.replace(/^\/api/, '') || '/'
  const dest = `${target}${stripped}${url.search}`
  return proxyRequest(event, dest)
})
