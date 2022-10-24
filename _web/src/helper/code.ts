const fullScreenKeys = {
  "F11": (cm: CodeMirror.Editor) => {
    cm.setOption("fullScreen", !cm.getOption("fullScreen"));
  },
  "Esc": (cm: CodeMirror.Editor) => {
    if (cm.getOption("fullScreen")) cm.setOption("fullScreen", false);
  },
  "Ctrl-/": (cm: CodeMirror.Editor) => {
    cm.toggleComment();
  }
};

export { fullScreenKeys };
