<script lang="ts">
  import { onDestroy, onMount, SvelteComponent } from "svelte";
  import update from "immutability-helper";
  import Router from "svelte-spa-router";
  import { storeData } from "@/store/store";
  import { push } from "svelte-spa-router";

  import Side from "./ui/Side.svelte";
  import Navbar from "./ui/Navbar.svelte";
  import Auths from "./pages/Auths.svelte";
  import Templates from "./pages/Templates.svelte";
  import Binds from "./pages/Binds.svelte";
  import Main from "./pages/Main.svelte";
  import { logout } from "@/helper/login";

  const routes = new Map<string | RegExp, typeof SvelteComponent>();
  routes.set(new RegExp("/auths(/(.*))*"), Auths);
  routes.set(new RegExp("/templates(/(.*))*"), Templates);
  routes.set(new RegExp("/binds(/(.*))*"), Binds);
  routes.set("*", Main);

  const sideLinks = ["auths", "templates", "binds"];

  let layout: HTMLElement;

  const select = (event: Event) => {
    if (event.target instanceof HTMLButtonElement) {
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
  });

  onDestroy(() => {
    layout.removeEventListener("click", select);
  });
</script>

<div class="layout h-full bg-gray-100" bind:this={layout}>
  <Navbar class="[grid-area:top]" />
  <Side class="[grid-area:sidebar]" links={sideLinks} />
  <div class="[grid-area:content] p-2 h-full">
    <Router {routes} />
  </div>
</div>

<style lang="scss">
  .layout {
    display: grid;
    grid-template-areas:
      "top top"
      "sidebar content"
      "sidebar content";
    grid-template-columns: 10rem 1fr;
    grid-template-rows: 3rem auto auto;
  }
</style>
