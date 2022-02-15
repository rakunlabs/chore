import { requestSender } from "@/helper/api";
import { tokenClear } from "@/helper/token";

const login = async (data:object) => {
  return requestSender(
    "/login",
    null,
    "POST",
    data,
    false,
    {
      timeout: 2000,
    }
  );
};

const renew = async () => {
  return requestSender(
    "/token/renew",
    null,
    "GET",
    null,
    true,
    {
      timeout: 2000,
    }
  );
};

const logout = () => {
  tokenClear();
};

export { login, renew, logout };
