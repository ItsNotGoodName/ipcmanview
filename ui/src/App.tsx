import { createSignal } from 'solid-js'

import { ExampleService } from "./core/client"

function App() {
  const [message, setMessage] = createSignal("Loading...");
  const api = new ExampleService(import.meta.env.VITE_BACKEND_URL, fetch)

  api.message().then(({ message }) => {
    setMessage(message.body + " " + message.time)
  }).catch((err) => {
    setMessage("Error: " + err)
  })

  return (
    <>
      <div>{message()}</div>
    </>
  )
}

export default App
