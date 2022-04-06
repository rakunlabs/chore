<script lang="ts">
  import Router, {
    push,
    ConditionsFailedEvent,
    querystring,
  } from "svelte-spa-router";
  import wrap from "svelte-spa-router/wrap";
  import { tokenCondition } from "@/helper/token";
  import Toast from "@/components/ui/Toast.svelte";
  import { pushRedirect } from "@/helper/push";

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
    if (event.detail.location == "/login") {
      pushRedirect($querystring);
      return;
    }

    push(`/login?back=${event.detail.location}`);
  };
</script>

<Toast />

<Router {routes} on:conditionsFailed={conditionsFailed} />
