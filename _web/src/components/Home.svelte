<script lang="ts">
  import { onDestroy, onMount, SvelteComponent } from "svelte";
  import update from "immutability-helper";
  import Router from "svelte-spa-router";
  import { push } from "svelte-spa-router";
  import { storeData } from "@/store/store";

  import Side from "@/components/ui/Side.svelte";
  import Navbar from "@/components/ui/Navbar.svelte";
  import Auths from "@/components/pages/Auths.svelte";
  import Templates from "@/components/pages/Templates.svelte";
  import Main from "@/components/pages/Main.svelte";
  import { logout } from "@/helper/login";
  import ControlFlow from "@/components/pages/ControlFlow.svelte";
  import Token from "@/components/pages/Token.svelte";
  import Head from "@/components/ui/Head.svelte";
  import Users from "@/components/pages/Users.svelte";
  import Send from "@/components/pages/Send.svelte";
  import Email from "@/components/pages/Email.svelte";
  import { isAdminToken } from "@/helper/token";

  // highlight operations
  import "@/helper/highlight";

  const routes = new Map<string | RegExp, typeof SvelteComponent>();
  routes.set(new RegExp("/send(/(.*))*"), Send);
  routes.set(new RegExp("/control(/(.*))*"), ControlFlow);
  routes.set(new RegExp("/auths(/(.*))*"), Auths);
  routes.set(new RegExp("/templates(/(.*))*"), Templates);
  routes.set(new RegExp("/token(/(.*))*"), Token);
  routes.set(new RegExp("/users(/(.*))*"), Users);
  routes.set(new RegExp("/email(/(.*))*"), Email);
  routes.set("*", Main);

  const sideLinks = [
    "send",
    "control",
    "auths",
    "templates",
    { settings: isAdminToken() ? ["token", "users", "email"] : ["token"] },
  ];

  let layout: HTMLElement;

  const select = (event: Event) => {
    if (event.target instanceof HTMLButtonElement) {
      const action = event.target.dataset["action"];
      if (action != "sidebar") {
        return;
      }

      const side = event.target.dataset["side"];

      if (side == "logout") {
        logout();
        push("/login");
        return;
      }

      // if same as before don't run it
      if (side != $storeData.sidebar) {
        storeData.update((v) =>
          update(v, {
            sidebar: { $set: side },
          })
        );
        push(`/${side}`);
      }
    }
  };

  onMount(() => {
    layout.addEventListener("click", select, false);
    storeData.update((v) =>
      update(v, {
        sidebar: { $set: "" },
      })
    );
  });

  onDestroy(() => {
    layout.removeEventListener("click", select);
  });
</script>

<div class="layout h-full bg-gray-100" bind:this={layout}>
  <Navbar class="[grid-area:top]" />
  <Side class="[grid-area:sidebar]" links={sideLinks} />
  <div class="[grid-area:content] h-full">
    <div class="grid h-full grid-rows-[auto_minmax(0,_1fr)]">
      <Head />
      <div class="p-2 h-full min-h-full overflow-y-auto">
        <Router {routes} />
      </div>
    </div>
  </div>
</div>

<style lang="scss">
  .layout {
    display: grid;
    grid-template-areas:
      "top top"
      "sidebar content";
    grid-template-columns: 10rem 1fr;
    grid-template-rows: 3rem minmax(0, 1fr);
  }
</style>
