import { AuthService } from "../core/client.gen"

const API = "http://localhost:8080"

const authService = new AuthService(API, fetch)

const USERNAME = "123"
const PASSWORD = "12345678"

try {
  const user = await authService.register({
    user: {
      email: "user@example.com",
      username: USERNAME,
      password: PASSWORD,
      passwordConfirm: PASSWORD,
    }
  })

  const { token } = await authService.login({
    usernameOrEmail: USERNAME,
    password: PASSWORD
  })

  console.log(token)
} catch (e) {
  console.log(e)
}
