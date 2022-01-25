const utf8ToB64 = ( s: string ) => {
  return btoa(unescape(encodeURIComponent( s )));
};

const b64ToUtf8 = ( s: string ) => {
  return decodeURIComponent(escape(atob( s )));
};

const formToObject = (form: HTMLFormElement) => {
  const formData = new FormData(form);
  const data: Record<string, any> = {};
  for (const field of formData) {
    const [key, value] = field;
    data[key] = value;
  }

  return data;
  // return JSON.stringify(data);
};

export { utf8ToB64, b64ToUtf8, formToObject };
