export default defineEventHandler(async (event) => {
  const body = await readBody<string>(event)

  if (body) {
    setCookie(event, "token", body, {
      path: "/",
      secure: true,
      sameSite: "strict",
      httpOnly: true,
      maxAge: 604800,
    })
  } else {
    deleteCookie(event, "token", {
      path: "/",
      secure: true,
      sameSite: "strict",
      httpOnly: true,
      maxAge: 604800,
    })
  }

  return ''
})

