import axios, { Method } from "axios";
import { trimLeft } from "./misc";

const requestSender = async (area: string, key: string, params: object, method: Method, data: any = undefined) => {
  const response = await axios({
    method: method,
    url: "./api/v1/kv/" + area,
    params: params,
    data: data,
  });

  return response;
};

const fixString = (area: string, key: string) => {
  if (area) {
    area = trimLeft(area);
  }
  if (key) {
    key = trimLeft(key);
  }

  return [area, key];
};

const getList = async (area: string, key: string) => {
  [area, key] = fixString(area, key);

  const params = {
    list: true,
    key: key,
  };

  return requestSender(area, key, params, "GET");
};

const getItem = async (area: string, key: string) => {
  [area, key] = fixString(area, key);

  const params = {
    key: key,
  };

  return requestSender(area, key, params, "GET");
};

const setItem = async (area: string, key: string, value: any) => {
  [area, key] = fixString(area, key);

  const params = {
    key: key,
  };

  return requestSender(area, key, params, "PUT", value);
};

const postItem = async (area: string, key: string, value: any) => {
  [area, key] = fixString(area, key);

  const params = {
    key: key,
  };

  return requestSender(area, key, params, "POST", value);
};


const deleteItem = async (area: string, key: string) => {
  [area, key] = fixString(area, key);

  const params = {
    key: key,
  };

  return requestSender(area, key, params, "DELETE");
};

export { getList, getItem, setItem, postItem, deleteItem };
