import { AuthService, WebrpcError } from "~/core/client.gen"

export default defineEventHandler(async (event) => {
  const runtimeConfig = useRuntimeConfig()
  const authService = new AuthService(runtimeConfig.public.apiBase, $fetch.raw)
  const session = getCookie(event, "session")

  if (!session) {
    return sendError(event, createError({ statusCode: 401 }))
  }

  try {
    const { token } = await authService.refresh({ session });
    return token
  } catch (e) {
    if (e instanceof WebrpcError) {
      return sendError(event, createError({ statusCode: e.status }))
    }

    return sendError(event, createError({ statusCode: 500 }))
  }
})

