import hljs from "highlight.js/lib/core";
import javascript from "highlight.js/lib/languages/javascript";
import xml from "highlight.js/lib/languages/xml";
import scss from "highlight.js/lib/languages/scss";
import bash from "highlight.js/lib/languages/bash";

// highlight in codemirror
import "codemirror/mode/javascript/javascript";
import "codemirror/mode/yaml/yaml";
import "codemirror/mode/shell/shell";
import "codemirror/addon/edit/matchbrackets";
import "codemirror/addon/edit/trailingspace";
import "codemirror/addon/selection/active-line";
import "codemirror/addon/display/fullscreen";
import "codemirror/addon/display/placeholder";

hljs.registerLanguage("javascript", javascript);
hljs.registerLanguage("http", xml);
hljs.registerLanguage("scss", scss);
hljs.registerLanguage("sh", bash);

export { hljs };
