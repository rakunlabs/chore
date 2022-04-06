const fullScreenKeys = {
  "F11": (cm: CodeMirror.Editor) => {
    cm.setOption("fullScreen", !cm.getOption("fullScreen"));
  },
  "Esc": (cm: CodeMirror.Editor) => {
    if (cm.getOption("fullScreen")) cm.setOption("fullScreen", false);
  },
};

export { fullScreenKeys };
