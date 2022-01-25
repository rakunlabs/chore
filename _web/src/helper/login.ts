import { requestSender } from "@/helper/api";
import { tokenClear } from "@/helper/token";

const login = async (data:object) => {
  return requestSender(
    "/login",
    null,
    "POST",
    data,
  );
};

const logout = () => {
  tokenClear();
};

export { login, logout };
