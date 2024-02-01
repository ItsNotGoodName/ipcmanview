/* @refresh reload */
import { render } from 'solid-js/web'

import App from "./App"

// https://github.com/GoogleChromeLabs/jsbi/issues/30#issuecomment-953187833
// @ts-ignore
BigInt.prototype.toJSON = function() { return this.toString() }

const root = document.getElementById('root')

render(() => <App />, root!)
