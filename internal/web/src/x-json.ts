import { LitElement, css, html, unsafeCSS } from 'lit'
import { customElement, property } from 'lit/decorators.js'

import hljs from 'highlight.js/lib/core';
import json from 'highlight.js/lib/languages/json';
import styles from "highlight.js/styles/tokyo-night-dark.css?inline"

hljs.registerLanguage('json', json);

@customElement('x-json')
export class JsonBlock extends LitElement {
  @property()
  json = '{}';

  render() {
    const template = document.createElement('code');
    template.innerHTML = hljs.highlight(this.json, { language: "json" }).value;
    template.className = "hljs"
    return html`
      <pre>${template}</pre>
    `
  }

  static get styles() {
    return [
      unsafeCSS(styles),
      css`
        pre { 
          margin: 0;
        }`
    ]
  }
}
