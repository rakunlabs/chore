import axios, { Method } from "axios";
import path from "path-browserify";
import { tokenGet } from "./token";

const requestSender = async (area: string, params: object, method: Method, data: any = undefined, useToken = false) => {
  let headers:Record<string, any> = {};

  if (useToken) {
    try {
      const [token] = tokenGet();
      headers = { Authorization: `Bearer ${token}` };
    } catch (error) {
      console.log(error);
    }
  }

  const response = await axios({
    method: method,
    url: path.join("./api/v1/", area),
    params: params,
    data: data,
    headers: headers,
  });

  return response;
};

// const fixString = (area: string, key: string) => {
//   if (area) {
//     area = trimLeft(area);
//   }
//   if (key) {
//     key = trimLeft(key);
//   }

//   return [area, key];
// };

// const getList = async (area: string, key: string) => {
//   [area, key] = fixString(area, key);

//   const params = {
//     list: true,
//     key: key,
//   };

//   return requestSender(area, key, params, "GET");
// };

// const getItem = async (area: string, key: string) => {
//   [area, key] = fixString(area, key);

//   const params = {
//     key: key,
//   };

//   return requestSender(area, key, params, "GET");
// };

// const setItem = async (area: string, key: string, value: any) => {
//   [area, key] = fixString(area, key);

//   const params = {
//     key: key,
//   };

//   return requestSender(area, key, params, "PUT", value);
// };

// const postItem = async (area: string, key: string, value: any) => {
//   [area, key] = fixString(area, key);

//   const params = {
//     key: key,
//   };

//   return requestSender(area, key, params, "POST", value);
// };


// const deleteItem = async (area: string, key: string) => {
//   [area, key] = fixString(area, key);

//   const params = {
//     key: key,
//   };

//   return requestSender(area, key, params, "DELETE");
// };

export { requestSender };
