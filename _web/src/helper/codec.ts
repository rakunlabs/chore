const formToObject = (form: HTMLFormElement) => {
  const formData = new FormData(form);
  const data: Record<string, any> = {};
  for (const field of formData) {
    const [key, value] = field;
    data[key] = value;
  }

  return data;
};

const formToObjectMulti = (form: HTMLFormElement) => {
  const formData = new FormData(form);
  const data: Record<string, any> = {};
  const remember: Record<string, any> = {};
  for (const field of formData) {
    const [key, value] = field;

    if (key.startsWith("headers-key")) {
      if (!data["headers"]) {
        data["headers"] = {} as Record<string, any>;
      }
      data["headers"][value.toString()] = "",
        remember[key.slice(11)] = value;
      continue;
    }

    if (key.startsWith("headers-value")) {
      const headerKey = remember[key.slice(13)];
      if (headerKey == null) {
        continue;
      }

      data["headers"][headerKey] = value;
      continue;
    }

    data[key] = value;
  }

  return data;
};

export { formToObject, formToObjectMulti };
