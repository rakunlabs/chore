import { addToast } from "@/store/toast";
import axios from "axios";
import type { CancelToken, Method } from "axios";
import path from "path-browserify";
import { tokenGet } from "./token";

const optionsDefault = {
  notTransformResponse: false,
  timeout: undefined as number,
  cancelToken: undefined as CancelToken,
  rawArea: false,
  noAlert: false,
};

const requestSender = async (area: string, params: object, method: Method, data: any = undefined, useToken = false, options: Partial<typeof optionsDefault> = optionsDefault) => {
  let headers: Record<string, any> = {};

  if (useToken) {
    try {
      const [token] = tokenGet();
      headers = { Authorization: `Bearer ${token["access_token"]}` };
    } catch (error) {
      console.log(error);
    }
  }

  try {
    const response = await axios({
      method: method,
      url: options.rawArea ? area : path.join("./api/v1/", area),
      params: params,
      data: data,
      headers: headers,
      timeout: options.timeout == null || options.timeout == undefined ? 2000 : options.timeout,
      transformResponse: options.notTransformResponse ? null : undefined,
      cancelToken: options.cancelToken,
    });

    return response;
  } catch (reason: unknown) {
    if (axios.isAxiosError(reason)) {
      const status = reason?.response?.status;
      if (!options.noAlert && status == 403) {
        addToast(reason?.response?.data?.error ?? reason?.message, "alert");
      }
    }

    throw reason;
  }
};

export { requestSender };
