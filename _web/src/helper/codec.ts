const utf8ToB64 = ( s: string ) => {
  return btoa(unescape(encodeURIComponent( s )));
};

const b64ToUtf8 = ( s: string ) => {
  return decodeURIComponent(escape(atob( s )));
};

export { utf8ToB64, b64ToUtf8 };
