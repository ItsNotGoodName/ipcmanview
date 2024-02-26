import { A, Route } from '@solidjs/router'

import { Home } from './Home'
import { Twirp } from './Twirp'
import { FormAction } from './FormAction'
import { loadHello } from './Twirp.data'
import { Cva } from './Cva'
import { Ui } from './Ui'
import { linkVariants } from '~/ui/Link'
import { Virtual } from './Virtual'

export default function() {
  return (
    <>
      <Route path="/home" component={Home} />
      <Route path="/twirp" component={Twirp} load={loadHello} />
      <Route path="/formaction" component={FormAction} />
      <Route path="/cva" component={Cva} />
      <Route path="/ui" component={Ui} />
      <Route path="/virtual" component={Virtual} />
      <Route path="/*" component={() =>
        <ul>
          <li><A class={linkVariants()} href='./home'>Home</A></li>
          <li><A class={linkVariants()} href='./twirp'>Twirp</A></li>
          <li><A class={linkVariants()} href='./formaction'>FormAction</A></li>
          <li><A class={linkVariants()} href='./cva'>Cva</A></li>
          <li><A class={linkVariants()} href='./ui'>Ui</A></li>
          <li><A class={linkVariants()} href='./virtual'>Virtual</A></li>
        </ul>
      } />
    </>
  )
}
