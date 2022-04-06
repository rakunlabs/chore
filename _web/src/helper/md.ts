import showdown from "showdown";

const converter = new showdown.Converter({
  // noHeaderId: true,
  simplifiedAutoLink: true,
  openLinksInNewWindow: true,
  prefixHeaderId: "/info-",
  rawPrefixHeaderId: true,
  extensions: [
    () => {
      const ancTpl = "$1<a class=\"anchor\" href=\"#$3\" aria-hidden=\"true\"><svg aria-hidden=\"true\" class=\"octicon octicon-link\" height=\"16\" version=\"1.1\" viewBox=\"0 0 16 16\" width=\"16\">"
      + "<path fill-rule=\"evenodd\" d=\"M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z\"></path>"
      +"</svg></a>$4";

      return [{
        type: "html",
        regex: /(<h([1-3]) id="([^"]+?)">)(.*<\/h\2>)/g,
        replace: ancTpl,
      }];
    },
  ],
});

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

export { converter, stringToHTML };
