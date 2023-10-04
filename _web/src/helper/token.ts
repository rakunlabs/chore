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
    let roles = claims?.roles as string[];

    if (roles == null) {
      roles = [];
    }

    // get all roles in the claims
    for (const resource in claims?.resource_access) {
      claims?.resource_access[resource]?.roles.forEach((role: string) => {
        roles.push(role);
      });
    }

    claims?.realm_access?.roles.forEach((role: string) => {
      roles.push(role);
    });


    return roles.includes("chore_admin");
  } catch (error) {
    console.error(error);
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

  return [data["token"], data["claims"], data["provider"]];
};

const tokenSet = (token: object, provider: string) => {
  const claims = jwtDecode(token["access_token"]);
  const data = JSON.stringify({
    provider,
    token,
    claims,
  });

  localStorage.setItem("token", data);
};

const tokenCondition = async () => {
  try {
    const [token] = tokenGet();
    await tokenCheck(token["access_token"]);
  } catch (error) {
    tokenClear();
    return false;
  }

  return true;
};

export { isAdminToken, tokenCheck, tokenClear, tokenGet, tokenSet, tokenCondition };
