import { requestSender } from "@/helper/api";

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

const tokenSet = (token: string, claims: object) => {
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

export { tokenCheck, tokenClear, tokenGet, tokenSet, tokenCondition };
