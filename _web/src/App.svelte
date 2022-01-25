<script lang="ts">
  import Router, { push, ConditionsFailedEvent } from "svelte-spa-router";
  import wrap from "svelte-spa-router/wrap";
  import { tokenCondition } from "@/helper/token";

  const routes = {
    "/login": wrap({
      asyncComponent: () => import("@/components/Login.svelte"),
      conditions: [
        async () => {
          const isTokenValid = await tokenCondition();
          return !isTokenValid;
        },
      ],
    }),
    "*": wrap({
      asyncComponent: () => import("@/components/Home.svelte"),
      conditions: [tokenCondition],
    }),
  };

  const conditionsFailed = (event: ConditionsFailedEvent) => {
    console.log(event);
    if (event.detail.location == "/login") {
      push("/");
      return;
    }

    push("/login");
  };
</script>

<Router {routes} on:conditionsFailed={conditionsFailed} />
