import showdown from "showdown";

import hljs from "highlight.js/lib/core";
import javascript from "highlight.js/lib/languages/javascript";
import xml from "highlight.js/lib/languages/xml";
import scss from "highlight.js/lib/languages/scss";
import "highlight.js/styles/monokai.css";

const converter = new showdown.Converter({
  noHeaderId: true,
  simplifiedAutoLink: true,
  openLinksInNewWindow: true,
});


hljs.registerLanguage("javascript", javascript);
hljs.registerLanguage("http", xml);
hljs.registerLanguage("scss", scss);

/**
 * Convert a template string into HTML DOM nodes
 * @param {string} str template string
 * @return {HTMLElement} template HTML
 */
const stringToHTML = (str: string) => {
  const parser = new DOMParser();
  const doc = parser.parseFromString(str, "text/html");
  return doc.body;
};

export { converter, hljs, stringToHTML };
