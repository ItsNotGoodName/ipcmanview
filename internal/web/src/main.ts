import "./index.css"

// ---------- htmx

import "./htmx.js"
import 'htmx.org/dist/ext/sse.js'

// ------------- Toastify

import Toastify from 'toastify-js'

document.body.addEventListener('htmx:afterRequest', function(evt: any) {
  if (evt.detail.failed) {
    const content = document.createElement("div")
    content.textContent = JSON.parse(evt.detail.xhr.responseText).message // evt.detail.xhr.responseText || evt.detail.xhr.statusText
    content.className = "flex-1"

    Toastify({
      text: JSON.parse(evt.detail.xhr.responseText).message, // evt.detail.xhr.responseText || evt.detail.xhr.statusText,
      node: content,
      duration: 3000,
      close: true,
      className: "alert alert-error flex flex-row",
      gravity: "bottom",
      position: "center",
      stopOnFocus: true,
    }).showToast();
  }
})

document.body.addEventListener("toast", function(evt: any) {
  const content = document.createElement("div")
  content.textContent = evt.detail.value
  content.className = "flex-1"

  Toastify({
    node: content,
    duration: 3000,
    close: true,
    className: "alert alert-success flex flex-row",
    gravity: "bottom",
    position: "center",
    stopOnFocus: true,
  }).showToast();
})

// ---------- Shoelace

// 24 Kb to format dates in local time
import '@shoelace-style/shoelace/dist/components/format-date/format-date.js';

// ---------- Lit

import "./x-json.ts"

// ---------- Checkbox

// @ts-ignore
window.tableCheckbox = function(ele: HTMLInputElement) {
  ele.
    closest("table")?.
    querySelectorAll(`th input[type="checkbox"]`).
    // @ts-ignore
    forEach((e: HTMLInputElement, _) => {
      if (!e.disabled)
        e.checked = ele.checked
    })
}
