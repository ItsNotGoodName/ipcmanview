// TODO: light and dark based on theme
import "highlight.js/styles/tokyo-night-dark.css"

import hljs from 'highlight.js/lib/core';
import json from 'highlight.js/lib/languages/json';

hljs.registerLanguage('json', json);

export default hljs 
