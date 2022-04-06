import { push } from "svelte-spa-router";

const pushRedirect = (querystring: string)=> {
  const params = new URLSearchParams(querystring);
  const backPath = params.get("back");

  if (backPath) {
    push(backPath);
  } else {
    push("/");
  }
};


export { pushRedirect };
