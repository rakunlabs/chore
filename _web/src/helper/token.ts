import { requestSender } from "@/helper/api";
import jwtDecode from "jwt-decode";

const tokenCheck = async (token: string) => {
  return requestSender(
    "/token/check",
    null,
    "POST",
    {
      token: token,
    },
  );
};

const tokenClear = () => {
  localStorage.removeItem("token");
};

const isAdminToken = () => {
  try {
    const [, claims] = tokenGet();
    const groups = claims["groups"] as string[];
    return groups ? groups.includes("admin") : false;
  } catch (error) {
    return false;
  }
};

const tokenGet = () => {
  const dataS = localStorage.getItem("token");
  if (dataS == null) {
    throw new Error("token not found");
  }

  const data = JSON.parse(dataS) as object;

  if (!data["token"]) {
    throw new Error("token not defined");
  }

  return [data["token"], data["claims"]];
};

const tokenSet = (token: string) => {
  const claims = jwtDecode(token);
  const data = JSON.stringify({
    token,
    claims,
  });

  localStorage.setItem("token", data);
};

const tokenCondition = async () => {
  try {
    const [token] = tokenGet();
    await tokenCheck(token);
  } catch (error) {
    tokenClear();
    return false;
  }

  return true;
};

export { isAdminToken, tokenCheck, tokenClear, tokenGet, tokenSet, tokenCondition };
