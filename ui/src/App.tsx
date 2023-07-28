import { createSignal } from 'solid-js'

import { ExampleService } from "./core/client"

function App() {
  const [message, setMessage] = createSignal("Loading...");
  const api = new ExampleService("http://localhost:3000", fetch)

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
