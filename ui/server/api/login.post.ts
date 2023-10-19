import { AuthService, LoginArgs, WebrpcError, } from "~/core/client.gen"
import { H3Event } from 'h3'

export type LoginRequest = Omit<LoginArgs, "clientId" | "ipAddress">

// TODO: add body types with Nuxt's h3 version is updated to v1.8.0
export default defineEventHandler(async (event) => {
  const body = await readBody<LoginRequest>(event)
  const runtimeConfig = useRuntimeConfig()
  const authService = new AuthService(runtimeConfig.public.apiBase, $fetch.raw)

  try {
    const clientId = getHeader(event, 'User-Agent') || ''
    const ipAddress = getRequestIP(event) || '0.0.0.0'
    const { res } = await authService.login({ ...body, clientId, ipAddress });

    setCookie(event, "session", res.session, {
      path: "/",
      secure: true,
      sameSite: "strict",
      httpOnly: true,
      expires: new Date(res.expiredAt)
    })

    return { token: res.token }
  } catch (e) {
    if (e instanceof WebrpcError) {
      return sendError(event, createError({ statusCode: e.status }))
    }

    return sendError(event, createError({ statusCode: 500 }))
  }
})

// TODO: delete when Nuxt's h3 version is updated to v1.8.0
// https://github.com/unjs/h3/blob/b688b22e11c95a810183e864053063e72815a944/src/utils/request.ts#L153C1-L180C2
function getRequestIP(
  event: H3Event,
  opts: {
    /**
     * Use the X-Forwarded-For HTTP header set by proxies.
     *
     * Note: Make sure that this header can be trusted (your application running behind a CDN or reverse proxy) before enabling.
     */
    xForwardedFor?: boolean;
  } = {},
): string | undefined {
  if (event.context.clientAddress) {
    return event.context.clientAddress;
  }

  if (opts.xForwardedFor) {
    const xForwardedFor = getRequestHeader(event, "x-forwarded-for")
      ?.split(",")
      ?.pop();
    if (xForwardedFor) {
      return xForwardedFor;
    }
  }

  if (event.node.req.socket.remoteAddress) {
    return event.node.req.socket.remoteAddress;
  }
}
