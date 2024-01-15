import { A, Route } from '@solidjs/router'
import { Home } from './Home'
import { Twirp } from './Twirp'
import { FormAction } from './FormAction'
import { loadHello } from './Twirp.data'
import { Cva } from './Cva'
import { Ui } from './Ui'

export function Debug() {
  return (
    <Route path="/debug">
      <Route path="/*" component={() =>
        <ul>
          <li><A href='./home'>Home</A></li>
          <li><A href='./twirp'>Twirp</A></li>
          <li><A href='./formaction'>FormAction</A></li>
          <li><A href='./cva'>Cva</A></li>
          <li><A href='./ui'>Ui</A></li>
        </ul>
      } />
      <Route path="/home" component={Home} />
      <Route path="/twirp" component={Twirp} load={loadHello} />
      <Route path="/formaction" component={FormAction} />
      <Route path="/cva" component={Cva} />
      <Route path="/ui" component={Ui} />
    </Route>
  )
}
