<script lang="ts">
  import Router, { push, querystring } from "svelte-spa-router";
  import type { ConditionsFailedEvent } from "svelte-spa-router";
  import wrap from "svelte-spa-router/wrap";
  import { tokenCondition } from "@/helper/token";
  import Toast from "@/components/ui/Toast.svelte";
  import { pushRedirect } from "@/helper/push";
  import { onMount } from "svelte";
  import { requestSender } from "@/helper/api";
  import { storeInfo } from "@/store/store";

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

  onMount(async () => {
    const response = await requestSender("/info", null, "GET");
    storeInfo.set(response.data);
  });
</script>

<Toast />

<Router {routes} on:conditionsFailed={conditionsFailed} />
