import { requestSender } from "@/helper/api";
import { tokenClear } from "@/helper/token";

const login = async (data: object, params: object | null) => {
  return requestSender(
    "/login",
    params,
    "POST",
    data,
    false,
  );
};

const renew = async (token: string, params: object | null) => {
  return requestSender(
    "/token/renew",
    params,
    "POST",
    { token },
    false,
  );
};

const logout = () => {
  tokenClear();
};

export { login, renew, logout };
