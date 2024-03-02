/* @refresh reload */
import { render } from 'solid-js/web'

import App from "./App"
import { BusProvider } from './providers/bus'
import { WSProvider } from './providers/ws'
import { ClientProvider } from './providers/client'

const root = document.getElementById('root')

render(() =>
  <BusProvider>
    <WSProvider>
      <ClientProvider>
        <App />
      </ClientProvider>
    </WSProvider>
  </BusProvider>
  , root!)
