import "./index.css";
import { Route, Router } from '@solidjs/router'
import { ClientProvider } from './providers/client'

import { useTheme } from "./ui/theme";
import { Layout } from './pages'
import { Debug } from './pages/debug'
import { NotFound } from './pages/404'
import { Home } from "./pages/Home";

function App() {
  useTheme()

  return (
    <ClientProvider>
      <Router base="/next" root={Layout}>
        <Debug />
        <Route path="/" component={Home} />
        <Route path="*404" component={NotFound} />
      </Router>
    </ClientProvider>
  )
}

export default App
