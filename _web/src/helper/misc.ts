const getBoolean = (value: string) => {
  if (value == null) {
    return false;
  }
  switch (value.toLocaleLowerCase()) {
  case "true":
    return true;
  default:
    return false;
  }
};

const getString = (value: boolean) => {
  return value ? "true" : "false";
};

const trimLeft = (k: string, v = "/") => k.replace(new RegExp(`^${v}`), "");
const trimRight = (k: string, v = "/") => k.replace(new RegExp(`$${v}`), "");

export { getBoolean, getString, trimLeft, trimRight };
