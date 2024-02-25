/* @refresh reload */
import { render } from 'solid-js/web'

import App from "./App"
import { BusProvider } from './providers/bus'
import { WSProvider } from './providers/ws'
import { ClientProvider } from './providers/client'

// https://github.com/GoogleChromeLabs/jsbi/issues/30#issuecomment-953187833
// @ts-expect-error
BigInt.prototype.toJSON = function() { return this.toString() }

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
