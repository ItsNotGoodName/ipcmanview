import { AuthService, WebrpcError } from "~/core/client.gen"

export default defineEventHandler(async (event) => {
  const session = getCookie(event, "session")
  if (!session) {
    return
  }

  const runtimeConfig = useRuntimeConfig()
  const authService = new AuthService(runtimeConfig.public.apiBase, $fetch.raw)

  try {
    await authService.logout({ session })
  } catch (e) {
    if (e instanceof WebrpcError) {
      return sendError(event, createError({ statusCode: e.status }))
    }

    return sendError(event, createError({ statusCode: 500 }))
  }

  deleteCookie(event, "session", {
    path: "/",
    secure: true,
    sameSite: "strict",
    httpOnly: true,
  })

  return ''
})

