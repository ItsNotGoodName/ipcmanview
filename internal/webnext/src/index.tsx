/* @refresh reload */
import { render } from 'solid-js/web'

import App from './App'
import { AuthProvider } from './providers/auth'
import { ClientProvider } from './providers/client'

const root = document.getElementById('root')

render(() => (
  <AuthProvider>
    <ClientProvider>
      <App />
    </ClientProvider>
  </AuthProvider>
), root!)
